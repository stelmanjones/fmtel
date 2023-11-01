package types

import (
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui/views"
	"github.com/stelmanjones/fmtel/units"
)

type App struct {
	Settings    Settings
	CarList     []cars.Car
	CurrentCar  cars.Car
	CurrentView views.View
}

type Settings struct {
	Temperature     units.Temperature
	EndpointAddress string
	UdpAddress      string
	EnableJSON      bool
	EnableSSE       bool
	RefreshRate     int
}
