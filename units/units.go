package units

import "time"

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
