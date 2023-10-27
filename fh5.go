package fmtel

import (
	"encoding/json"
	"time"

	"github.com/stelmanjones/fmtel/units"
)

type Horizon5Packet struct {
	// ----------------------- SLED
	// --------------------------------
	IsRaceOn    int32  // = 1 when race is on, = 0 when in menus/race stopped
	TimestampMs uint32 // can overflow to 0 eventually

	// ~ General Information

	DriveTrainType      uint8 // Between 100 (slowest car) and 999 (fastest car) inclusive
	EngineCylinders     uint8 // Number of cylinders in the engine
	CarPerformanceIndex uint8 // Between 100 (slowest car) and 999 (fastest car) inclusive
	CarClass            uint8 // Between 0 (D -- worst cars) and 7 (X class -- best cars) inclusive
	CarOrdinal          uint8 // Unique ID of the car make/model

	// ~ Engine Information

	EngineMaxRPM     float32
	EngineIdleRPM    float32
	CurrentEngineRPM float32

	// ~ Wheel information

	WheelRotationSpeedFrontLeft  float32
	WheelRotationSpeedFrontRight float32
	WheelRotationSpeedRearLeft   float32
	WheelRotationSpeedRearRight  float32

	WheelOnRumbleStripFrontLeft  uint8
	WheelOnRumbleStripFrontRight uint8
	WheelOnRumbleStripRearLeft   uint8
	WheelOnRumbleStripRearRight  uint8

	WheelInPuddleDepthFrontLeft  float32
	WheelInPuddleDepthFrontRight float32
	WheelInPuddleDepthRearLeft   float32
	WheelInPuddleDepthRearRight  float32

	// ~ Tire information

	TireSlipRotationFrontLeft  float32
	TireSlipRotationFrontRight float32
	TireSlipRotationRearLeft   float32
	TireSlipRotationRearRight  float32

	TireSlipAngleFrontLeft  float32
	TireSlipAngleFrontRight float32
	TireSlipAngleRearLeft   float32
	TireSlipAngleRearRight  float32

	TireCombinedSlipFrontLeft  float32
	TireCombinedSlipFrontRight float32
	TireCombinedSlipRearLeft   float32
	TireCombinedSlipRearRight  float32

	TireTempFrontLeft  float32
	TireTempFrontRight float32
	TireTempRearLeft   float32
	TireTempRearRight  float32

	// ~ Suspension Information

	NormalizedSuspensionTravelFrontLeft  float32
	NormalizedSuspensionTravelFrontRight float32
	NormalizedSuspensionTravelRearLeft   float32
	NormalizedSuspensionTravelRearRight  float32

	SuspensionTravelMetersFrontLeft  float32
	SuspensionTravelMetersFrontRight float32
	SuspensionTravelMetersRearLeft   float32
	SuspensionTravelMetersRearRight  float32

	// ~ Spatial information

	PositionX float32
	PositionY float32
	PositionZ float32

	AccelerationX float32
	AccelerationY float32
	AccelerationZ float32

	VelocityX float32
	VelocityY float32
	VelocityZ float32

	AngularVelocityX float32
	AngularVelocityY float32
	AngularVelocityZ float32

	Yaw   float32
	Pitch float32
	Roll  float32

	// ~ Force feedback information

	SurfaceRumbleFrontLeft  float32
	SurfaceRumbleFrontRight float32
	SurfaceRumbleRearLeft   float32
	SurfaceRumbleRearRight  float32

	// ----------------------- DASH
	// --------------------------------

	// ~ Literal dashboard information

	Speed  float32 // meters/second
	Power  float32 // watts
	Torque float32 // newton meter

	Boost            float32
	Fuel             float32
	DistanceTraveled float32

	Throttle  uint8
	Brake     uint8
	Clutch    uint8
	Handbrake uint8
	Gear      uint8
	Steer     int8

	// ~ Lap information

	LapNumber       uint16
	BestLap         float32
	LastLap         float32
	CurrentLap      float32
	CurrentRaceTime float32
	RacePosition    uint8

	// ~ Game data

	NormalizedDrivingLine       uint8
	NormalizedAIBrakeDifference uint8
}

// Returns current racetime as a formatted string. "03:42.583"
func (f *Horizon5Packet) FmtCurrentRaceTime() (t string) {
	return units.Timespan(time.Duration(f.CurrentRaceTime * float32(time.Second))).Format("04:05.000")
}

// Returns current laptime as a formatted string. "03:42.583"
func (f *Horizon5Packet) FmtCurrentLap() (t string) {
	return units.Timespan(time.Duration(f.CurrentLap * float32(time.Second))).Format("04:05.000")
}

// Returns last laptime as a formatted string. "03:42.583"
func (f *Horizon5Packet) FmtLastLap() (t string) {
	return units.Timespan(time.Duration(f.LastLap * float32(time.Second))).Format("04:05.000")
}

// Returns best laptime as a formatted string. "03:42.583"
func (f *Horizon5Packet) FmtBestLap() (t string) {
	return units.Timespan(time.Duration(f.BestLap * float32(time.Second))).Format("04:05.000")
}

// Returns current suspension travel in Meters.
func (f *Horizon5Packet) SuspensionTravelMeters() (s *SuspensionTravel) {
	s = &SuspensionTravel{
		Normalized: false,
		FrontLeft:  f.SuspensionTravelMetersFrontLeft,
		FrontRight: f.SuspensionTravelMetersFrontRight,
		RearLeft:   f.SuspensionTravelMetersRearLeft,
		RearRight:  f.SuspensionTravelMetersRearRight,
	}
	return
}

// Returns current suspension travel as a value betweeen 0(no travel) and 1.0(max travel)
func (f *Horizon5Packet) NormalizedSuspensionTravel() (s *SuspensionTravel) {
	s = &SuspensionTravel{
		Normalized: true,
		FrontLeft:  f.NormalizedSuspensionTravelFrontLeft,
		FrontRight: f.NormalizedSuspensionTravelFrontRight,
		RearLeft:   f.NormalizedSuspensionTravelRearLeft,
		RearRight:  f.NormalizedSuspensionTravelRearRight,
	}
	return
}

// Returns the current coordinates of the car
func (f *Horizon5Packet) CarPosition() (p *Position) {
	p = &Position{}
	p.X = f.PositionX
	p.Y = f.PositionY
	p.Z = f.PositionZ
	return
}

// Returns true if game is paused or not in a race
func (f *Horizon5Packet) IsPaused() bool {
	switch f.IsRaceOn {
	case 1:
		return true
	default:
		return false
	}
}

// Returns current engine power output in horsepower
func (m *Horizon5Packet) HorsePower() uint {
	return uint(m.Power * 0.00134102)
}

// Returns current engine power output in kilowatts
func (m *Horizon5Packet) Kilowatts() uint {
	return uint(m.Power / 1000)
}

// Returns current speed in mph
func (m *Horizon5Packet) MilesPerHour() uint {
	return uint(m.Speed * 2.2369362921)
}

// Returns current speed in kmph
func (m *Horizon5Packet) KmPerHour() uint {
	return uint(m.Speed * 3.6)
}

// Returns current engine torque in ft/lbs
func (m *Horizon5Packet) FootPounds() uint {
	return uint(float64(m.Torque) / 1.356)
}

// Returns the current cars drivetrain type as a label ( FWD , RWD , AWD )
// If the type cannot be parsed it returns "-"
func (m *Horizon5Packet) FmtDrivetrainType() string {
	switch m.DriveTrainType {
	case 0:
		{
			switch m.IsRaceOn {
			case 1:
				return "FWD"
			default:
				return "-"
			}
		}
	case 1:
		return "RWD"
	case 2:
		return "AWD"
	default:
		return "-"
	}
}

// Returns current cars class as a formatted label ("D","C","B","A","S","R","P","X") or "-"
func (x *Horizon5Packet) FmtCarClass() string {
	switch x.CarClass - 1 {
	case 0:
		if x.IsRaceOn == 1 {
			return "D"
		} else {
			return "-"
		}
	case 1:
		return "C"
	case 2:
		return "B"
	case 3:
		return "A"
	case 4:
		return "S"
	case 5:
		return "R"
	case 6:
		return "P"
	case 7:
		return "X"
	default:
		return "-"
	}
}

func (m *Horizon5Packet) ToJson() ([]byte, error) {
	return json.Marshal(m)
}
