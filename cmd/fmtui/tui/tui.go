package tui

import (
	"fmt"
	"time"

	//"math"

	//"github.com/pterm/pterm/putils"
	//flag "github.com/spf13/pflag"
	"github.com/charmbracelet/log"
	"github.com/pterm/pterm"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cmd/fmtui/graph"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"
	"github.com/stelmanjones/fmtel/units"
)

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
			temps = *packet.TireTempsInCelsius()
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

	inputGraph := graph.RenderInputGraph(app.GraphData, *packet, app.GraphDataPoints)
	boost := func() float32 {
		if packet.Boost < 1 {
			return 0.0
		} else {
			return packet.Boost
		}
	}()

	var currentTime time.Duration = time.Duration(packet.CurrentRaceTime * float32(time.Second))
	var currentLapTime time.Duration = time.Duration(packet.CurrentLap * float32(time.Second))
	var bestLapTime time.Duration = time.Duration(packet.BestLap * float32(time.Second))
	var lastLapTime time.Duration = time.Duration(packet.LastLap * float32(time.Second))

	stats, err := pterm.DefaultTable.WithLeftAlignment().WithData(pterm.TableData{
		{
			"Drivetrain Type: ", packet.ParsedDrivetrainType(),
			"Car Class:", packet.ParsedCarClass(),
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
		{"Laptime:", units.Timespan(currentLapTime).Format("04:05.000")},
		{"Last Lap:", units.Timespan(lastLapTime).Format("04:05.000")},
		{"Best Lap:", units.Timespan(bestLapTime).Format("04:05.000")},
		{"Current Racetime:", units.Timespan(currentTime).Format("15:04:05.00")},
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
				Sprintf("\n\nFMTEL | Version: 0.1.0 \n\n%s %s %s\n\n", pterm.FgWhite.Sprint(currentCar.Maker), currentCar.Model, pterm.FgDarkGray.ToStyle().Sprintf("(%d)", currentCar.Year)))
	tires := WheelTempWidget(packet, &app.Settings)
	layout, err := pterm.DefaultPanel.WithPadding(4).WithPanels(pterm.Panels{
		{{Data: title}},
		{{Data: pterm.DefaultBox.WithTitle("Race Info").WithBoxStyle(pterm.FgLightBlue.ToStyle()).Sprint(lapStats)}, {Data: tires}},
		{{Data: pterm.Sprintf("%s", stats)}},
		{{Data: pterm.Sprintf("\n\n%s", inputGraph)}},
	}).Srender()
	if err != nil {
		log.Error(err)
	}
	return layout
}
