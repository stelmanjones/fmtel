package types

import (
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/units"
)

type App struct {
	Settings        Settings
	CarList         []cars.Car
	GraphData       [][]float64
	CurrentCar      cars.Car
	GraphDataPoints int
}

type Settings struct {
	Temperature units.Temperature
	UdpAddress  string
}

