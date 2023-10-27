package tui

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/pterm/pterm"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cmd/fmtui/pedals"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"
	"github.com/stelmanjones/fmtel/units"
)

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

	final := pterm.DefaultBox.WithTitleTopLeft().WithTitle("Pedals").WithBoxStyle(pterm.FgGreen.ToStyle()).Sprint(pedals)
	return final
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

	final := pterm.DefaultBox.WithBoxStyle(pterm.FgWhite.ToStyle()).WithTitle("Tire Temps").Sprint(template)
	return final
}

func Render(packet *fmtel.ForzaPacket, app *types.App) string {
	currentCar := app.CurrentCar

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
			"PI:", fmt.Sprintf("%3d", packet.CarPerformanceIndex),
		},

		{
			"RPM:", fmt.Sprintf("%5.f rpm", packet.CurrentEngineRpm),
			"Horsepower: ", fmt.Sprintf("%5d hp", packet.HorsePower()),
			"Kilowatts: ", fmt.Sprintf("%5d kw", packet.Kilowatts()),
		},

		{
			"Gear:", fmt.Sprintf("%5d", packet.Gear),
			"Max RPM:", fmt.Sprintf("%5.f rpm", packet.EngineMaxRpm),
			"Car ID:", fmt.Sprintf("%5d", packet.CarOrdinal),
		},

		{
			"Speed (kmh):", fmt.Sprintf("%5d km/h", packet.KmPerHour()),
			"Idle RPM:", fmt.Sprintf("%5.f rpm", packet.EngineIdleRpm),
			"Torque (nm):", fmt.Sprintf("%5d nm", uint(packet.Torque)),
		},

		{
			"Speed (mph):", fmt.Sprintf("%5d mph", packet.MilesPerHour()),
			"Boost", fmt.Sprintf("%5.f psi", boost),
			"Torque (ft/lb):", fmt.Sprintf("%5d ft/lb", packet.FootPounds()),
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
	title := pterm.DefaultCenter.
		WithCenterEachLineSeparately(true).
		Sprint(
			pterm.
				FgGreen.
				ToStyle().
				Add(*pterm.
					Bold.
					ToStyle()).
				Sprintf("\n\nFMTEL | Version: 0.1.1 \n\n%s %s %s\n\n", pterm.FgWhite.Sprint(currentCar.Maker), currentCar.Model, pterm.FgDarkGray.ToStyle().Sprintf("(%d)", currentCar.Year)))
	tires := WheelTempWidget(packet, &app.Settings)
	layout, err := pterm.DefaultPanel.WithPadding(4).WithBottomPadding(2).WithPanels(pterm.Panels{
		{{Data: title}},
		{{Data: pterm.DefaultBox.WithTitle("Race Info").WithBoxStyle(pterm.FgLightBlue.ToStyle()).Sprint(lapStats)}, {Data: tires}},
		{{Data: pterm.Sprintf("%s", stats)}},
		{{Data: pedals}},
	}).Srender()
	if err != nil {
		log.Error(err)
	}
	return layout
}
