package graph

import (
	"github.com/guptarohit/asciigraph"
	"github.com/pterm/pterm"
	"github.com/stelmanjones/fmtel"
)

// 0 = Throttle
// 1 = Brake
// 2 = Clutch
func InputPercentage(val uint8) float64 {
	return (float64(val) / 255) * 100
}

func RenderInputGraph(data [][]float64, packet fmtel.ForzaPacket, points int) string {
	vals := []float64{InputPercentage(packet.Accel), InputPercentage(packet.Brake), InputPercentage(packet.Clutch)}
	for i := 0; i < 2; i++ {
		data[i] = append(data[i], vals[i])
		if points > 0 && len(data[i]) > points {
			data[i] = data[i][len(data[i])-points:]
		}
	}
	return asciigraph.PlotMany(data, asciigraph.AxisColor(asciigraph.DimGray),
		asciigraph.SeriesColors(
			asciigraph.GreenYellow,
			asciigraph.Red,
			asciigraph.Yellow),
		asciigraph.LowerBound(0),
		asciigraph.UpperBound(100),
		asciigraph.Height(5),
		asciigraph.Precision(0),
		asciigraph.Width(50),
		asciigraph.Offset(5),
		asciigraph.LabelColor(asciigraph.DarkGray),
		asciigraph.Caption(pterm.Sprintf("      %s | %s | %s %s",
			pterm.FgGreen.Sprint("Throttle"),
			pterm.FgLightRed.Sprint("Brake"),
			pterm.FgYellow.Sprint("Clutch"),
			pterm.FgDarkGray.Sprint("(%)"))))
}
