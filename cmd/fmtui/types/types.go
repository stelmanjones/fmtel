package types

import (
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/units"
)

type App struct {
	Settings   Settings
	CarList    []cars.Car
	CurrentCar cars.Car
}

type Settings struct {
	Temperature units.Temperature
	UdpAddress  string
}
