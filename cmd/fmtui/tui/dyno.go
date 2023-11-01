package tui

import (
	"sync"

	"github.com/pterm/pterm"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui/views"
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
// Gear 0 = Reverse
type Gears [11]Speed

func (g *Gears) Update(p *fmtel.ForzaPacket) {
	if p.Gear > 10 {
		return
	}
	gear := p.Gear
	if g[gear].Kmh < p.KmPerHour() {
		g[gear].Kmh = p.KmPerHour()
		g[gear].Mph = p.MilesPerHour()
	}
}

func (g *Gears) Reset() {
	for i := 0; i <= len(g)-1; i++ {
		g[i].Kmh = 0
		g[i].Mph = 0
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
	d.TopSpeedPerGear.Reset()
	d.TopSpeed = NewSpeed()
	d.MaxFootPounds = 0
	d.MaxNewtonMeters = 0
	d.MaxHorsePower = 0
	d.MaxKiloWatts = 0
}

func RenderDynoView(d *DynoView, a *cars.Car) (string, error) {
	return pterm.DefaultPanel.WithBottomPadding(2).WithPanels(pterm.Panels{
		{{Data: TitleBar()}},
		{{Data: CarInfo(a)}},
		{{
			Data: pterm.DefaultBasicText.Sprintf(`
    
    gear 1: %d\n
    gear 2: %d\n
    gear 3: %d\n
    gear 4: %d\n
    gear 5: %d\n
    gear 6: %d\n
    gear 7: %d\n
    gear 8: %d\n
    gear 9: %d\n
    gear 10: %d\n
    rev: %d\n
    `, d.TopSpeedPerGear[1],
				d.TopSpeedPerGear[2],
				d.TopSpeedPerGear[3],
				d.TopSpeedPerGear[4],
				d.TopSpeedPerGear[5],
				d.TopSpeedPerGear[6],
				d.TopSpeedPerGear[7],
				d.TopSpeedPerGear[8],
				d.TopSpeedPerGear[9],
				d.TopSpeedPerGear[10],
				d.TopSpeedPerGear[0]),
		}},
		{{Data: StatusBar(views.Dyno)}},
	}).WithSameColumnWidth(true).Srender()
}
