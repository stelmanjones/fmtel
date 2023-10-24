package pkg

import (
	"encoding/json"
)

// TODO: Improve documentation.
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
type TireTemperatures struct {
	FrontLeft  float32
	FrontRight float32
	RearLeft   float32
	RearRight  float32
}

type PedalInputs struct {
	Clutch   uint
	Brake    uint
	Throttle uint
}

func (m *ForzaPacket) IsPaused() bool {
	switch m.IsRaceOn {
	case 1:
		return true
	default:
		return false
	}
}

func (m *ForzaPacket) HorsePower() uint {
	return uint(m.Power * 0.00134102)
}

func (m *ForzaPacket) Kilowatts() uint {
	return uint(m.Power / 1000)
}

func (m *ForzaPacket) MilesPerHour() uint {
	return uint(m.Speed * 2.2369362921)
}

func (m *ForzaPacket) KmPerHour() uint {
	return uint(m.Speed * 3.6)
}

func (m *ForzaPacket) FootPounds() uint {
	return uint(float64(m.Torque) / 1.356)
}

func (m *ForzaPacket) TireTempsInCelsius() *TireTemperatures {
	b := TireTemperatures{
		(m.TireTempFrontLeft - 32) * 5 / 9,
		(m.TireTempFrontRight - 32) * 5 / 9,
		(m.TireTempRearLeft - 32) * 5 / 9,
		(m.TireTempRearRight - 32) * 5 / 9,
	}
	return &b
}

// Returns the current cars drivetrain type as a label ( FWD , RWD , AWD ).
// If the type cannot be parsed it returns "-".
func (m *ForzaPacket) ParsedDrivetrainType() string {
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

func (x *ForzaPacket) ParsedCarClass() string {
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
