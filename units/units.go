package units

import (
	"strconv"
	"strings"
	"time"
)

type (
	Units       uint
	Temperature string
	Drivetrain  string
)

const (
	FWD     Drivetrain = "FWD"
	RWD     Drivetrain = "RWD"
	AWD     Drivetrain = "AWD"
	UNKNOWN Drivetrain = "Unknown"
)

const (
	CELSIUS    Temperature = "celsius"
	FAHRENHEIT Temperature = "fahrenheit"
)

const (
	METRIC Units = iota
	IMPERIAL
)

func TempFromString(t string) Temperature {
	switch t {
	case "fahrenheit":
		return FAHRENHEIT
	default:
		return CELSIUS
	}
}

// Helper type for formatting of laptimes.
type Timespan time.Duration

func (t Timespan) Format(format string) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
}

func PadItoa(val int, length int) string {
	b := strings.Builder{}
	parsed := strconv.Itoa(val)
	diff := length - len(parsed)
	b.WriteString(parsed)
	for i := 1; i <= diff; i++ {
		b.WriteRune('\u2000')
	}
	return b.String()
}

func PadUInt(val uint, length int) string {
	b := strings.Builder{}
	parsed := strconv.FormatUint(uint64(val), 10)
	diff := length - len(parsed)
	b.WriteString(parsed)
	for i := 1; i <= diff; i++ {
		b.WriteRune('\u2000')
	}
	return b.String()
}
