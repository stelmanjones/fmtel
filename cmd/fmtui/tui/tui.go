package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/pterm/pterm"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/cmd/fmtui/pedals"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui/views"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"
	"github.com/stelmanjones/fmtel/internal"
	"github.com/stelmanjones/fmtel/units"
)

func TitleBar() (t string) {
	t = pterm.DefaultCenter.
		WithCenterEachLineSeparately(true).
		Sprint(
			pterm.
				FgGreen.
				ToStyle().
				Add(*pterm.
					Bold.
					ToStyle()).
				Sprintf("\n\nFMTEL | %s", internal.VERSION))
	return
}

func CarInfo(currentCar *cars.Car) (c string) {
	c = pterm.DefaultCenter.
		WithCenterEachLineSeparately(true).
		Sprint(
			pterm.
				FgGreen.
				ToStyle().
				Add(*pterm.
					Bold.
					ToStyle()).
				Sprint(pterm.FgWhite.Sprint(currentCar.Maker), " ", currentCar.Model, " ", pterm.FgDarkGray.ToStyle().Sprintf("(%s)", units.PadItoa(int(currentCar.Year), 4))))
	return
}

func StatusBar(v views.View) string {
	switch v {
	default:
		{
			return pterm.DefaultBasicText.
				WithStyle(pterm.FgDarkGray.ToStyle()).
				Sprint("ctrl+c/q/escape: quit • t: switch °C/°F • d: toggle dyno")
		}
	case views.Dyno:
		{
			return pterm.DefaultBasicText.
				WithStyle(pterm.FgDarkGray.ToStyle()).
				Sprint("ctrl+c/q/escape: quit • d: toggle dyno • r: reset dyno")
		}

	}
}

func PedalWidget(packet *fmtel.ForzaPacket) string {
	pedals, err := pedals.DefaultPedalInputBar.WithHeight(20).WithWidth(30).WithBars(pterm.Bars{
		pterm.Bar{
			Label: "Throttle",
			Value: int(packet.Accel),
			Style: pterm.FgGreen.ToStyle(),
		},
		pterm.Bar{
			Label: "Brake",
			Value: int(packet.Brake),
			Style: pterm.FgRed.ToStyle(),
		},
		pterm.Bar{
			Label: "Clutch",
			Value: int(packet.Clutch),
			Style: pterm.FgYellow.ToStyle(),
		},
	}).Srender()
	if err != nil {
		log.Error(err)
	}

	return pterm.DefaultBox.WithTitleTopLeft().WithTitle("Pedals").WithBoxStyle(pterm.FgGreen.ToStyle()).Sprint(pedals)
}

func WheelTempWidget(packet *fmtel.ForzaPacket, settings *types.Settings) string {
	var temps fmtel.TireTemperatures

	switch settings.Temperature {
	case "fahrenheit":
		temps = fmtel.TireTemperatures{
			FrontLeft:  packet.TireTempFrontLeft,
			FrontRight: packet.TireTempFrontRight,
			RearLeft:   packet.TireTempRearLeft,
			RearRight:  packet.TireTempRearRight,
		}
	default:
		{
			temps = *packet.TireTempsCelsius()
		}
	}

	var template string
	switch settings.Temperature {

	case units.CELSIUS:
		template = pterm.Sprintf("\n        F  \n%3.f°C █   █ %3.f°C \n\n%3.f°C █   █ %3.f°C \n        R\n", temps.FrontLeft, temps.FrontRight, temps.RearLeft, temps.RearRight)
	case units.FAHRENHEIT:
		template = pterm.Sprintf("\n        F  \n%3.f°F █   █ %3.f°F \n\n%3.f°F █   █ %3.f°F \n        R\n", temps.FrontLeft, temps.FrontRight, temps.RearLeft, temps.RearRight)
	}

	return pterm.DefaultBox.WithBoxStyle(pterm.FgWhite.ToStyle()).WithTitle("Tire Temps").Sprint(template)
}

// WARN: Refactor ASAP. MEMORY HOG
func Render(packet *fmtel.ForzaPacket, app *types.App) string {
	boost := func() float32 {
		if packet.Boost <= 0 {
			return 0.0
		} else {
			return packet.Boost
		}
	}()

	pedals := PedalWidget(packet)

	stats, err := pterm.DefaultTable.WithLeftAlignment().WithBoxed(true).WithData(pterm.TableData{
		{
			"Drivetrain Type: ", packet.FmtDrivetrainType(),
			"Car Class:", packet.FmtCarClass(),
			"PI:", units.PadItoa(int(packet.CarPerformanceIndex), 4),
		},

		{
			"RPM:", strconv.FormatFloat(float64(packet.CurrentEngineRpm), 'f', 0, 64) + " rpm",
			"Horsepower: ", units.PadItoa(int(packet.HorsePower()), 5) + " hp",
			"Kilowatts: ", units.PadItoa(int(packet.KiloWatts()), 4) + " kw",
		},

		{
			"Gear:", strconv.Itoa(int(packet.Gear)),
			"Max RPM:", units.PadItoa(int(packet.EngineMaxRpm), 5) + " rpm",
			"Car ID:", units.PadItoa(int(packet.CarOrdinal), 5),
		},

		{
			"Speed (kmh):", units.PadItoa(int(packet.KmPerHour()), 3) + " km/h",
			"Idle RPM:", units.PadItoa(int(packet.EngineIdleRpm), 5) + " rpm",
			"Torque (nm):", units.PadUInt(uint(packet.Torque), 4) + " nm",
		},

		{
			"Speed (mph):", units.PadItoa(int(packet.MilesPerHour()), 3) + " mph",
			"Boost", fmt.Sprintf("%s psi", strconv.FormatFloat(float64(boost), 'f', 3, 64)),
			"Torque (ft/lb):", units.PadItoa(int(packet.FootPounds()), 4) + " ft/lb",
		},
	}).Srender()
	if err != nil {
		log.Error(err)
	}

	lapStats, err := pterm.DefaultTable.WithLeftAlignment().WithData(pterm.TableData{
		{"Postition:", fmt.Sprintf("%2d", packet.RacePosition)},
		{"Lap: ", fmt.Sprintf("%2d", packet.LapNumber)},
		{"Laptime:", packet.FmtCurrentLap()},
		{"Last Lap:", packet.FmtLastLap()},
		{"Best Lap:", packet.FmtBestLap()},
		{"Current Racetime:", packet.FmtCurrentRaceTime()},
	}).Srender()
	if err != nil {
		log.Error(err)
	}

	tires := WheelTempWidget(packet, &app.Settings)
	layout, err := pterm.DefaultPanel.WithBottomPadding(2).WithPanels(pterm.Panels{
		{{Data: TitleBar()}},
		{{Data: CarInfo(&app.CurrentCar)}},
		{{Data: pterm.DefaultBox.WithTitle("Race Info").WithBoxStyle(pterm.FgLightBlue.ToStyle()).Sprint(lapStats)}, {Data: tires}},
		{{Data: pterm.Sprint(stats)}},
		{{Data: pedals}},
		{{Data: StatusBar(views.Home)}},
	}).Srender()
	if err != nil {
		log.Error(err)
	}
	return layout
}
