package tui

import (
	"sync"

	"github.com/stelmanjones/fmtel"
)

type Speed struct {
	Kmh uint
	Mph uint
}

func NewSpeed() Speed {
	return Speed{
		Kmh: 0,
		Mph: 0,
	}
}

// Helper struct that keeps track of all gears.
// Gear 11(index 10) = Reverse
type Gears [11]Speed

func (g *Gears) Update(p *fmtel.ForzaPacket) {
	gear := func() uint8 {
		if p.Gear <= 0 {
			return 1
		} else {
			return p.Gear
		}
	}()
	if g[gear].Kmh < p.KmPerHour() {
		g[gear].Kmh = p.KmPerHour()
		g[gear].Mph = p.MilesPerHour()
	}
}

func NewGears() Gears {
	return Gears{
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
		NewSpeed(),
	}
}

type DynoView struct {
	mu              *sync.RWMutex
	TopSpeedPerGear Gears
	TopSpeed        Speed
	MaxNewtonMeters float32
	MaxFootPounds   float32
	MaxHorsePower   float32
	MaxKiloWatts    float32
}

func NewDynoView() *DynoView {
	return &DynoView{
		&sync.RWMutex{},
		NewGears(),
		NewSpeed(),
		0,
		0,
		0,
		0,
	}
}

func (d *DynoView) Update(p *fmtel.ForzaPacket) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.TopSpeedPerGear.Update(p)

	if d.TopSpeed.Kmh < p.KmPerHour() {
		d.TopSpeed.Kmh = p.KmPerHour()
		d.TopSpeed.Mph = p.MilesPerHour()
	}
	if d.MaxNewtonMeters < p.Torque {
		d.MaxNewtonMeters = p.Torque
		d.MaxFootPounds = float32(p.FootPounds())
	}
	if d.MaxHorsePower < float32(p.HorsePower()) {
		d.MaxHorsePower = float32(p.HorsePower())
		d.MaxKiloWatts = float32(p.KiloWatts())
	}

}

func (d *DynoView) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.TopSpeedPerGear = NewGears()
	d.TopSpeed = NewSpeed()
	d.MaxFootPounds = 0
	d.MaxNewtonMeters=	0
		d.MaxHorsePower = 0
	d.MaxKiloWatts = 0

}