package units

import "time"

type Units uint
type Temperature string
type Drivetrain string

const (
	FWD Drivetrain = "FWD"
	RWD Drivetrain = "RWD"
	AWD Drivetrain = "AWD"
	UNKNOWN Drivetrain = "Unknown"

)

const (
	CELSIUS Temperature = "celsius"
	FAHRENHEIT Temperature = "fahrenheit"
)

const (
	METRIC Units = iota
	IMPERIAL
)

func TempFromString(t string) Temperature {
	switch t {
		case "fahrenheit": return FAHRENHEIT
		default: return CELSIUS
	}
}
type Timespan time.Duration
func (t Timespan) Format(format string) string {
    return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
}