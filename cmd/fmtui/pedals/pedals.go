package pedals

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/gookit/color"
	"github.com/mattn/go-runewidth"
	"github.com/pterm/pterm"
)

var (
	// Output completely disables output from pterm if set to false. Can be used in CLI application quiet mode.
	Output = true

	// PrintDebugMessages sets if messages printed by the DebugPrinter should be printed.
	PrintDebugMessages = false

	// RawOutput is set to true if pterm.DisableStyling() was called.
	// The variable indicates that PTerm will not add additional styling to text.
	// Use pterm.DisableStyling() or pterm.EnableStyling() to change this variable.
	// Changing this variable directly, will disable or enable the output of colored text.
	RawOutput = false
)

var (
	// ErrTerminalSizeNotDetectable - the terminal size can not be detected and the fallback values are used.
	ErrTerminalSizeNotDetectable = errors.New("terminal size could not be detected - using fallback value")

	// ErrHexCodeIsInvalid - the given HEX code is invalid.
	ErrHexCodeIsInvalid = errors.New("hex code is not valid")
)

// FallbackTerminalWidth is the value used for GetTerminalWidth, if the actual width can not be detected
// You can override that value if necessary.
var FallbackTerminalWidth = 80

// FallbackTerminalHeight is the value used for GetTerminalHeight, if the actual height can not be detected
// You can override that value if necessary.
var FallbackTerminalHeight = 10

// forcedTerminalWidth, when set along with forcedTerminalHeight, forces the terminal width value.
var forcedTerminalWidth int = 0

// forcedTerminalHeight, when set along with forcedTerminalWidth, forces the terminal height value.
var forcedTerminalHeight int = 0

func RemoveColorFromString(a ...interface{}) string {
	return color.ClearCode(fmt.Sprint(a...))
}

func GetTerminalWidth() int {
	if forcedTerminalWidth > 0 {
		return forcedTerminalWidth
	}
	width, _, _ := GetTerminalSize()
	return width
}

// GetTerminalHeight returns the terminal height of the active terminal.
func GetTerminalHeight() int {
	if forcedTerminalHeight > 0 {
		return forcedTerminalHeight
	}
	_, height, _ := GetTerminalSize()
	return height
}

func GetTerminalSize() (width, height int, err error) {
	if forcedTerminalWidth > 0 && forcedTerminalHeight > 0 {
		return forcedTerminalWidth, forcedTerminalHeight, nil
	}
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if w <= 0 {
		w = FallbackTerminalWidth
	}
	if h <= 0 {
		h = FallbackTerminalHeight
	}
	if err != nil {
		err = ErrTerminalSizeNotDetectable
	}
	return w, h, err
}

func GetStringMaxWidth(s string) int {
	var max int
	ss := strings.Split(s, "\n")
	for _, s2 := range ss {
		s2WithoutColor := color.ClearCode(s2)
		if runewidth.StringWidth(s2WithoutColor) > max {
			max = runewidth.StringWidth(s2WithoutColor)
		}
	}
	return max
}

// setForcedTerminalSize turns off terminal size autodetection. Usuful for unified tests.
func SetForcedTerminalSize(width int, height int) {
	forcedTerminalWidth = width
	forcedTerminalHeight = height
	RecalculateTerminalSize()
}

func RecalculateTerminalSize() {
	// keep in sync with DefaultBarChart
	DefaultPedalInputBar.Width = GetTerminalWidth() * 2 / 3
	DefaultPedalInputBar.Height = GetTerminalHeight() * 2 / 3
	pterm.DefaultParagraph.MaxWidth = GetTerminalWidth()
}

func MapRangeToRange(fromMin, fromMax, toMin, toMax, current float32) int {
	if fromMax-fromMin == 0 {
		return 0
	}
	return int(toMin + ((toMax-toMin)/(fromMax-fromMin))*(current-fromMin))
}

func WithBoolean(b []bool) bool {
	if len(b) == 0 {
		b = append(b, true)
	}
	return b[0]
}

// PedalInputBar is used to print live pedal data. (Throttle,Brake,Clutch,Handbrake)
type PedalInputBar struct {
	Writer                 io.Writer
	VerticalBarCharacter   string
	HorizontalBarCharacter string
	Bars                   pterm.Bars
	Height                 int
	Width                  int
	Horizontal             bool
	ShowValue              bool
}

// DefaultBarChart is the default PedalInputBar.
var DefaultPedalInputBar = PedalInputBar{
	Horizontal:             true,
	VerticalBarCharacter:   "██",
	HorizontalBarCharacter: "█",
	// keep in sync with RecalculateTerminalSize()
	Height: GetTerminalHeight() * 2 / 3,
	Width:  GetTerminalWidth() * 2 / 3,
}

// WithBars returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithBars(bars pterm.Bars) *PedalInputBar {
	p.Bars = bars
	return &p
}

// WithVerticalBarCharacter returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithVerticalBarCharacter(char string) *PedalInputBar {
	p.VerticalBarCharacter = char
	return &p
}

// WithHorizontalBarCharacter returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithHorizontalBarCharacter(char string) *PedalInputBar {
	p.HorizontalBarCharacter = char
	return &p
}

// WithHorizontal returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithHorizontal(b ...bool) *PedalInputBar {
	b2 := WithBoolean(b)
	p.Horizontal = b2
	return &p
}

// WithHeight returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithHeight(value int) *PedalInputBar {
	p.Height = value
	return &p
}

// WithWidth returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithWidth(value int) *PedalInputBar {
	p.Width = value
	return &p
}

// WithShowValue returns a new PedalInputBar with a specific option.
func (p PedalInputBar) WithShowValue(b ...bool) *PedalInputBar {
	p.ShowValue = WithBoolean(b)
	return &p
}

// WithWriter sets the custom Writer.
func (p PedalInputBar) WithWriter(writer io.Writer) *PedalInputBar {
	p.Writer = writer
	return &p
}

func (p PedalInputBar) getRawOutput() string {
	var ret string

	for _, bar := range p.Bars {
		ret += pterm.Sprintfln("%s: %d", bar.Label, bar.Value)
	}

	return ret
}

// Srender renders the BarChart as a string.
func (p PedalInputBar) Srender() (string, error) {
	//_maxAbsValue := 100

	abs := func(value int) int {
		if value < 0 {
			return -value
		}

		return value
	}
	// =================================== VERTICAL BARS RENDERER ======================================================

	type renderParams struct {
		bar                     pterm.Bar
		indent                  string
		repeatCount             int
		positiveChartPartHeight int
		negativeChartPartHeight int
		positiveChartPartWidth  int
		negativeChartPartWidth  int
		showValue               bool
		moveUp                  bool
		moveRight               bool
	}

	renderPositiveVerticalBar := func(renderedBarRef *string, rParams renderParams) {
		if rParams.showValue {
			*renderedBarRef += fmt.Sprint(rParams.indent + strconv.Itoa(rParams.bar.Value) + rParams.indent + "\n")
		}

		for i := rParams.positiveChartPartHeight; i > 0; i-- {
			if i > rParams.repeatCount {
				*renderedBarRef += rParams.indent + "  " + rParams.indent + " \n"
			} else {
				*renderedBarRef += rParams.indent + rParams.bar.Style.Sprint(p.VerticalBarCharacter) + rParams.indent + " \n"
			}
		}

		// Used when we draw diagram with both POSITIVE and NEGATIVE values.
		// In such case we separately draw top and bottom half of chart.
		// And we need MOVE UP positive part to top part of chart,
		// technically by adding empty pillars with height == height of chart's bottom part.
		if rParams.moveUp {
			for i := 0; i <= rParams.negativeChartPartHeight; i++ {
				*renderedBarRef += rParams.indent + "  " + rParams.indent + " \n"
			}
		}
	}

	renderNegativeVerticalBar := func(renderedBarRef *string, rParams renderParams) {
		for i := 0; i > -rParams.negativeChartPartHeight; i-- {
			if i > rParams.repeatCount {
				*renderedBarRef += rParams.indent + rParams.bar.Style.Sprint(p.VerticalBarCharacter) + rParams.indent + " \n"
			} else {
				*renderedBarRef += rParams.indent + "  " + rParams.indent + " \n"
			}
		}

		if rParams.showValue {
			*renderedBarRef += fmt.Sprint(rParams.indent + strconv.Itoa(rParams.bar.Value) + rParams.indent + "\n")
		}
	}

	// =================================== HORIZONTAL BARS RENDERER ====================================================
	renderPositiveHorizontalBar := func(renderedBarRef *string, rParams renderParams) {
		if rParams.moveRight {
			for i := 0; i < rParams.negativeChartPartWidth; i++ {
				*renderedBarRef += " "
			}
		}

		for i := 0; i < rParams.positiveChartPartWidth; i++ {
			if i < rParams.repeatCount {
				*renderedBarRef += rParams.bar.Style.Sprint(p.HorizontalBarCharacter)
			} else {
				*renderedBarRef += " "
			}
		}

		if rParams.showValue {
			// For positive horizontal bars we add one more space before adding value,
			// so they will be well aligned with negative values, which have "-" sign before them
			*renderedBarRef += " "

			*renderedBarRef += " " + strconv.Itoa(rParams.bar.Value)
		}
	}

	renderNegativeHorizontalBar := func(renderedBarRef *string, rParams renderParams) {
		for i := -rParams.negativeChartPartWidth; i < 0; i++ {
			if i < rParams.repeatCount {
				*renderedBarRef += " "
			} else {
				*renderedBarRef += rParams.bar.Style.Sprint(p.HorizontalBarCharacter)
			}
		}

		// In order to print values well-aligned (in case when we have both - positive and negative part of chart),
		// we should insert an indent with width == width of positive chart part
		if rParams.positiveChartPartWidth > 0 {
			for i := 0; i < rParams.positiveChartPartWidth; i++ {
				*renderedBarRef += " "
			}
		}

		if rParams.showValue {
			/*
				This is in order to achieve this effect:
				 0
				-15
				 0
				-19

				INSTEAD OF THIS:

				0
				-15
				0
				-19
			*/
			if rParams.repeatCount == 0 {
				*renderedBarRef += " "
			}

			*renderedBarRef += " " + strconv.Itoa(rParams.bar.Value)
		}
	}
	// =================================================================================================================

	if RawOutput {
		return p.getRawOutput(), nil
	}
	for i, bar := range p.Bars {
		if bar.Style == nil {
			p.Bars[i].Style = &pterm.ThemeDefault.BarStyle
		}

		if bar.LabelStyle == nil {
			p.Bars[i].LabelStyle = &pterm.ThemeDefault.BarLabelStyle
		}

		p.Bars[i].Label = p.Bars[i].LabelStyle.Sprint(bar.Label)
	}

	var ret string

	var maxLabelHeight int
	maxBarValue := 255
	minBarValue := 0
	maxAbsBarValue := 255
	var rParams renderParams

	for _, bar := range p.Bars {
		if bar.Value > maxBarValue {
			maxBarValue = bar.Value
		}
		if bar.Value < minBarValue {
			minBarValue = bar.Value
		}
		labelHeight := len(strings.Split(bar.Label, "\n"))
		if labelHeight > maxLabelHeight {
			maxLabelHeight = labelHeight
		}
	}

	maxAbsBarValue = 100

	if p.Horizontal {
		panels := pterm.Panels{[]pterm.Panel{{}, {}}}

		rParams.showValue = p.ShowValue
		rParams.positiveChartPartWidth = p.Width
		rParams.negativeChartPartWidth = p.Width

		// If chart will consist of two parts - positive and negative - we should recalculate max bars WIDTH in LEFT and RIGHT parts
		if minBarValue < 0 && maxBarValue > 0 {
			rParams.positiveChartPartWidth = abs(MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Width)/2, float32(p.Width)/2, float32(maxBarValue)))
			rParams.negativeChartPartWidth = abs(MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Width)/2, float32(p.Width)/2, float32(minBarValue)))
		}

		for _, bar := range p.Bars {
			rParams.bar = bar
			panels[0][0].Data += "\n" + bar.Label
			panels[0][1].Data += "\n"

			if minBarValue >= 0 {
				// As we don't have negative values, draw only positive (right) part of the chart:
				rParams.repeatCount = MapRangeToRange(0, float32(maxAbsBarValue), 0, float32(p.Width), float32(bar.Value))
				rParams.moveRight = false

				renderPositiveHorizontalBar(&panels[0][1].Data, rParams)
			} else if maxBarValue <= 0 {
				// As we have only negative values, draw only negative (left) part of the chart:
				rParams.repeatCount = MapRangeToRange(-float32(maxAbsBarValue), 0, -float32(p.Width), 0, float32(bar.Value))
				rParams.positiveChartPartWidth = 0

				renderNegativeHorizontalBar(&panels[0][1].Data, rParams)
			} else {
				// We have positive and negative values, so draw both (left+right) parts of the chart:
				rParams.repeatCount = MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Width)/2, float32(p.Width)/2, float32(bar.Value))

				if bar.Value >= 0 {
					rParams.moveRight = true

					renderPositiveHorizontalBar(&panels[0][1].Data, rParams)
				}

				if bar.Value < 0 {
					renderNegativeHorizontalBar(&panels[0][1].Data, rParams)
				}
			}
		}
		ret, _ = pterm.DefaultPanel.WithPanels(panels).Srender()
		return ret, nil
	} else {
		renderedBars := make([]string, len(p.Bars))

		rParams.showValue = p.ShowValue
		rParams.positiveChartPartHeight = p.Height
		rParams.negativeChartPartHeight = p.Height

		// If chart will consist of two parts - positive and negative - we should recalculate max bars height in top and bottom parts
		if minBarValue < 0 && maxBarValue > 0 {
			rParams.positiveChartPartHeight = abs(MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Height)/2, float32(p.Height)/2, float32(maxBarValue)))
			rParams.negativeChartPartHeight = abs(MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Height)/2, float32(p.Height)/2, float32(minBarValue)))
		}

		for i, bar := range p.Bars {
			var renderedBar string
			rParams.bar = bar
			rParams.indent = strings.Repeat(" ", GetStringMaxWidth(RemoveColorFromString(bar.Label))/2)

			if minBarValue >= 0 {
				// As we don't have negative values, draw only positive (top) part of the chart:
				rParams.repeatCount = MapRangeToRange(0, float32(maxAbsBarValue), 0, float32(p.Height), float32(bar.Value))
				rParams.moveUp = false // Don't MOVE UP as we have ONLY positive part of chart.

				renderPositiveVerticalBar(&renderedBar, rParams)
			} else if maxBarValue <= 0 {
				// As we have only negative values, draw only negative (bottom) part of the chart:
				rParams.repeatCount = MapRangeToRange(-float32(maxAbsBarValue), 0, -float32(p.Height), 0, float32(bar.Value))

				renderNegativeVerticalBar(&renderedBar, rParams)
			} else {
				// We have positive and negative values, so draw both (top+bottom) parts of the chart:
				rParams.repeatCount = MapRangeToRange(-float32(maxAbsBarValue), float32(maxAbsBarValue), -float32(p.Height)/2, float32(p.Height)/2, float32(bar.Value))

				if bar.Value >= 0 {
					rParams.moveUp = true // MOVE UP positive part, because we have both positive and negative parts of chart.

					renderPositiveVerticalBar(&renderedBar, rParams)
				}

				if bar.Value < 0 {
					renderNegativeVerticalBar(&renderedBar, rParams)
				}
			}

			labelHeight := len(strings.Split(bar.Label, "\n"))
			renderedBars[i] = renderedBar + bar.Label + strings.Repeat("\n", maxLabelHeight-labelHeight) + " "
		}

		var maxBarHeight int

		for _, bar := range renderedBars {
			totalBarHeight := len(strings.Split(bar, "\n"))
			if totalBarHeight > maxBarHeight {
				maxBarHeight = totalBarHeight
			}
		}

		for i, bar := range renderedBars {
			totalBarHeight := len(strings.Split(bar, "\n"))
			if totalBarHeight < maxBarHeight {
				renderedBars[i] = strings.Repeat("\n", maxBarHeight-totalBarHeight) + renderedBars[i]
			}
		}

		for i := 0; i <= maxBarHeight; i++ {
			for _, barString := range renderedBars {
				var barLine string
				letterLines := strings.Split(barString, "\n")
				maxBarWidth := GetStringMaxWidth(RemoveColorFromString(barString))
				if len(letterLines) > i {
					barLine = letterLines[i]
				}
				letterLineLength := runewidth.StringWidth(RemoveColorFromString(barLine))
				if letterLineLength < maxBarWidth {
					barLine += strings.Repeat(" ", maxBarWidth-letterLineLength)
				}
				ret += barLine
			}
			ret += "\n"
		}
	}

	return ret, nil
}

// Render prints the Template to the terminal.
func (p PedalInputBar) Render() error {
	s, _ := p.Srender()
	pterm.Fprintln(p.Writer, s)

	return nil
}
