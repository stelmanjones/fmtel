package fmtel

import (
	"encoding/json"
	"time"

	"github.com/stelmanjones/fmtel/units"
)

type ForzaPacket struct {
	// = 1 when race is on. = 0 when in menus/race stopped
	IsRaceOn int32
	// Can overflow to 0 eventually
	TimestampMS      uint32
	EngineMaxRpm     float32
	EngineIdleRpm    float32
	CurrentEngineRpm float32

	// In the car's local space; X = right, Y = up, Z = forward
	AccelerationX float32
	AccelerationY float32
	AccelerationZ float32

	// In the car's local space; X = right, Y = up, Z = forward
	VelocityX float32
	VelocityY float32
	VelocityZ float32

	// In the car's local space; X = pitch, Y = yaw, Z = roll
	AngularVelocityX float32
	AngularVelocityY float32
	AngularVelocityZ float32
	Yaw              float32
	Pitch            float32
	Roll             float32

	// Suspension travel normalized: 0.0f = max stretch; 1.0 = max compression
	NormalizedSuspensionTravelFrontLeft  float32
	NormalizedSuspensionTravelFrontRight float32
	NormalizedSuspensionTravelRearLeft   float32
	NormalizedSuspensionTravelRearRight  float32

	// Tire normalized slip ratio, = 0 means 100% grip and |ratio| > 1.0 means loss of grip.
	TireSlipRatioFrontLeft  float32
	TireSlipRatioFrontRight float32
	TireSlipRatioRearLeft   float32
	TireSlipRatioRearRight  float32

	// Wheels rotation speed radians/sec.
	WheelRotationSpeedFrontLeft  float32
	WheelRotationSpeedFrontRight float32
	WheelRotationSpeedRearLeft   float32
	WheelRotationSpeedRearRight  float32

	// = 1 when wheel is on rumble strip, = 0 when off.
	WheelOnRumbleStripFrontLeft  int32
	WheelOnRumbleStripFrontRight int32
	WheelOnRumbleStripRearLeft   int32
	WheelOnRumbleStripRearRight  int32

	// = from 0 to 1, where 1 is the deepest puddle
	WheelInPuddleDepthFrontLeft  float32
	WheelInPuddleDepthFrontRight float32
	WheelInPuddleDepthRearLeft   float32
	WheelInPuddleDepthRearRight  float32

	// Non-dimensional surface rumble values passed to controller force feedback
	SurfaceRumbleFrontLeft  float32
	SurfaceRumbleFrontRight float32
	SurfaceRumbleRearLeft   float32
	SurfaceRumbleRearRight  float32

	// Tire normalized slip angle, = 0 means 100% grip and |angle| > 1.0 means loss of grip.
	TireSlipAngleFrontLeft  float32
	TireSlipAngleFrontRight float32
	TireSlipAngleRearLeft   float32
	TireSlipAngleRearRight  float32

	// Tire normalized combined slip, = 0 means 100% grip and |slip| > 1.0 means loss of grip.
	TireCombinedSlipFrontLeft  float32
	TireCombinedSlipFrontRight float32
	TireCombinedSlipRearLeft   float32
	TireCombinedSlipRearRight  float32

	// Actual suspension travel in meters
	SuspensionTravelMetersFrontLeft  float32
	SuspensionTravelMetersFrontRight float32
	SuspensionTravelMetersRearLeft   float32
	SuspensionTravelMetersRearRight  float32

	// Car ID
	CarOrdinal int32

	// Between 0 (D -- worst cars) and 7 (X class -- best cars) inclusive
	CarClass int32

	// Between 100 (worst car) and 999 (best car) inclusive
	CarPerformanceIndex int32

	// 0 = FWD, 1 = RWD, 2 = AWD
	DrivetrainType int32

	// Number of cylinders in the engine
	NumCylinders int32

	PositionX float32
	PositionY float32
	PositionZ float32
	// Speed in meters per second.
	Speed float32
	// Power in kilowatts.
	Power float32
	// Torque in newtonmeters.
	Torque float32

	// Tire temperatures in fahrenheit.
	TireTempFrontLeft  float32
	TireTempFrontRight float32
	TireTempRearLeft   float32
	TireTempRearRight  float32
	// Boost in psi.
	Boost            float32
	Fuel             float32
	DistanceTraveled float32
	BestLap          float32
	LastLap          float32
	CurrentLap       float32
	CurrentRaceTime  float32
	LapNumber        uint16
	RacePosition     uint8
	// Between 0(none) - 255(full).
	Accel uint8
	// Between 0(none) - 255(full).
	Brake uint8
	// Between 0(none) - 255(full).
	Clutch uint8
	// Between 0(none) - 255(full).
	HandBrake uint8
	Gear      uint8
	// Between -255(left) - 255(right).
	Steer int8

	NormalizedDrivingLine       int8
	NormalizedAIBrakeDifference int8
	TireWearFrontLeft           float32
	TireWearFrontRight          float32
	TireWearRearLeft            float32
	TireWearRearRight           float32
	// Unique track ID.
	TrackOrdinal int32
}

type (
	TireWear struct {
		FrontLeft  float32
		FrontRight float32
		RearLeft   float32
		RearRight  float32
	}

	TireTemperatures struct {
		FrontLeft  float32
		FrontRight float32
		RearLeft   float32
		RearRight  float32
	}

	PedalInputs struct {
		Clutch   uint
		Brake    uint
		Throttle uint
	}

	SuspensionTravel struct {
		Normalized bool
		FrontLeft  float32
		FrontRight float32
		RearLeft   float32
		RearRight  float32
	}

	Position struct {
		X float32
		Y float32
		Z float32
	}
)

// Returns current racetime as a formatted string. "03:42.583"
func (f *ForzaPacket) FmtCurrentRaceTime() (t string) {
	return units.Timespan(time.Duration(f.CurrentRaceTime * float32(time.Second))).Format("04:05.000")
}

// Returns current laptime as a formatted string. "03:42.583"
func (f *ForzaPacket) FmtCurrentLap() (t string) {
	return units.Timespan(time.Duration(f.CurrentLap * float32(time.Second))).Format("04:05.000")
}

// Returns last laptime as a formatted string. "03:42.583"
func (f *ForzaPacket) FmtLastLap() (t string) {
	return units.Timespan(time.Duration(f.LastLap * float32(time.Second))).Format("04:05.000")
}

// Returns best laptime as a formatted string. "03:42.583"
func (f *ForzaPacket) FmtBestLap() (t string) {
	return units.Timespan(time.Duration(f.BestLap * float32(time.Second))).Format("04:05.000")
}

// Returns current suspension travel in Meters.
func (f *ForzaPacket) SuspensionTravelMeters() (s *SuspensionTravel) {
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
func (f *ForzaPacket) NormalizedSuspensionTravel() (s *SuspensionTravel) {
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
func (f *ForzaPacket) CarPosition() (p *Position) {
	p = &Position{}
	p.X = f.PositionX
	p.Y = f.PositionY
	p.Z = f.PositionZ
	return
}

// Returns true if game is paused or not in a race
func (f *ForzaPacket) IsPaused() bool {
	switch f.IsRaceOn {
	case 1:
		return true
	default:
		return false
	}
}

// Returns current engine power output in horsepower
func (m *ForzaPacket) HorsePower() uint {
	return uint(m.Power * 0.00134102)
}

// Returns current engine power output in kilowatts
func (m *ForzaPacket) Kilowatts() uint {
	return uint(m.Power / 1000)
}

// Returns current speed in mph
func (m *ForzaPacket) MilesPerHour() uint {
	return uint(m.Speed * 2.2369362921)
}

// Returns current speed in kmph
func (m *ForzaPacket) KmPerHour() uint {
	return uint(m.Speed * 3.6)
}

// Returns current engine torque in ft/lbs
func (m *ForzaPacket) FootPounds() uint {
	return uint(float64(m.Torque) / 1.356)
}

// Returns current tire wear. Between 1.0(no wear) and 0(max wear)
func (f *ForzaPacket) TireWear() (t *TireWear) {
	t = &TireWear{
		FrontLeft:  f.TireWearFrontLeft,
		FrontRight: f.TireWearFrontRight,
		RearLeft:   f.TireWearRearLeft,
		RearRight:  f.TireWearRearRight,
	}
	return
}

// Returns current tire temperatures in Celsius
func (m *ForzaPacket) TireTempsCelsius() *TireTemperatures {
	b := TireTemperatures{
		(m.TireTempFrontLeft - 32) * 5 / 9,
		(m.TireTempFrontRight - 32) * 5 / 9,
		(m.TireTempRearLeft - 32) * 5 / 9,
		(m.TireTempRearRight - 32) * 5 / 9,
	}
	return &b
}

// Returns the current cars drivetrain type as a label ( FWD , RWD , AWD )
// If the type cannot be parsed it returns "-"
func (m *ForzaPacket) FmtDrivetrainType() string {
	switch m.DrivetrainType {
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
func (x *ForzaPacket) FmtCarClass() string {
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

func (m *ForzaPacket) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

var DefaultForzaPacket = ForzaPacket{
	IsRaceOn:                             0,
	TimestampMS:                          0,
	EngineMaxRpm:                         0,
	EngineIdleRpm:                        0,
	CurrentEngineRpm:                     0,
	AccelerationX:                        0,
	AccelerationY:                        0,
	AccelerationZ:                        0,
	VelocityX:                            0,
	VelocityY:                            0,
	VelocityZ:                            0,
	AngularVelocityX:                     0,
	AngularVelocityY:                     0,
	AngularVelocityZ:                     0,
	Yaw:                                  0,
	Pitch:                                0,
	Roll:                                 0,
	NormalizedSuspensionTravelFrontLeft:  0,
	NormalizedSuspensionTravelFrontRight: 0,
	NormalizedSuspensionTravelRearLeft:   0,
	NormalizedSuspensionTravelRearRight:  0,
	TireSlipRatioFrontLeft:               0,
	TireSlipRatioFrontRight:              0,
	TireSlipRatioRearLeft:                0,
	TireSlipRatioRearRight:               0,
	WheelRotationSpeedFrontLeft:          0,
	WheelRotationSpeedFrontRight:         0,
	WheelRotationSpeedRearLeft:           0,
	WheelRotationSpeedRearRight:          0,
	WheelOnRumbleStripFrontLeft:          0,
	WheelOnRumbleStripFrontRight:         0,
	WheelOnRumbleStripRearLeft:           0,
	WheelOnRumbleStripRearRight:          0,
	WheelInPuddleDepthFrontLeft:          0,
	WheelInPuddleDepthFrontRight:         0,
	WheelInPuddleDepthRearLeft:           0,
	WheelInPuddleDepthRearRight:          0,
	SurfaceRumbleFrontLeft:               0,
	SurfaceRumbleFrontRight:              0,
	SurfaceRumbleRearLeft:                0,
	SurfaceRumbleRearRight:               0,
	TireSlipAngleFrontLeft:               0,
	TireSlipAngleFrontRight:              0,
	TireSlipAngleRearLeft:                0,
	TireSlipAngleRearRight:               0,
	TireCombinedSlipFrontLeft:            0,
	TireCombinedSlipFrontRight:           0,
	TireCombinedSlipRearLeft:             0,
	TireCombinedSlipRearRight:            0,
	SuspensionTravelMetersFrontLeft:      0,
	SuspensionTravelMetersFrontRight:     0,
	SuspensionTravelMetersRearLeft:       0,
	SuspensionTravelMetersRearRight:      0,
	CarOrdinal:                           0,
	CarClass:                             0,
	CarPerformanceIndex:                  0,
	DrivetrainType:                       0,
	NumCylinders:                         0,
	PositionX:                            0,
	PositionY:                            0,
	PositionZ:                            0,
	Speed:                                0,
	Power:                                0,
	Torque:                               0,
	TireTempFrontLeft:                    0,
	TireTempFrontRight:                   0,
	TireTempRearLeft:                     0,
	TireTempRearRight:                    0,
	Boost:                                0,
	Fuel:                                 0,
	DistanceTraveled:                     0,
	BestLap:                              0,
	LastLap:                              0,
	CurrentLap:                           0,
	CurrentRaceTime:                      0,
	LapNumber:                            0,
	RacePosition:                         0,
	Accel:                                0,
	Brake:                                0,
	Clutch:                               0,
	HandBrake:                            0,
	Gear:                                 0,
	Steer:                                0,
	NormalizedDrivingLine:                0,
	NormalizedAIBrakeDifference:          0,
	TireWearFrontLeft:                    0,
	TireWearFrontRight:                   0,
	TireWearRearLeft:                     0,
	TireWearRearRight:                    0,
	TrackOrdinal:                         0,
}
