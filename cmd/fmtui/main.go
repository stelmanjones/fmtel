package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/muesli/termenv"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard/keys"
	flag "github.com/spf13/pflag"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cmd/fmtui/input"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"

	"github.com/charmbracelet/log"
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/server"
	"github.com/stelmanjones/fmtel/units"
	// "github.com/pterm/pterm"
)

var Pack = fmtel.DefaultForzaPacket

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

var (
	temp        string
	udpAddress  string
	enableJson  bool
	jsonAddress string
	upgrader    = websocket.Upgrader{} // use default options
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, _, err := c.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		data, err := Pack.ToJson()
		if err != nil {
			log.Error("Serialization error:", err)
			break
		}
		err = c.WriteMessage(mt, data)
		if err != nil {
			log.Error("write:", err)
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func responder(w http.ResponseWriter, r *http.Request) {
	data, err := Pack.ToJson()
	if err != nil {
		log.Error(err)
	}

	switch r.Method {
	case "GET":
		enableCors(&w)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Not supported.")
	}
}

func serveHTTP() {
	http.HandleFunc("/forza/ws", wsHandler)
	http.HandleFunc("/forza", responder)

	log.Debugf("JSON Telemetry Server started at http://localhost%s", ":9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	flag.StringVar(&temp, "temp", "celsius", "Set temperature unit.")
	flag.StringVar(&udpAddress, "udp-addr", ":7777", "Set UDP connection address.")
	flag.StringVar(&jsonAddress, "json-addr", ":9999", "Set JSON server address.")
	flag.BoolVar(&enableJson, "json", false, "Enable JSON server.")
	flag.Lookup("json").NoOptDefVal = "true"
	flag.Parse()

	out := termenv.DefaultOutput()
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(out)
	if err != nil {
		panic(err)
	}
	out.AltScreen()
	cursor.Hide()

	defer cursor.Show()

	settings := types.Settings{
		Temperature: units.TempFromString(temp),
		UdpAddress:  udpAddress,
	}
	carList, err := cars.ReadCarList("cars.json")
	if err != nil {
		log.Error(err)
	}
	app := types.App{
		CurrentCar:      cars.DefaultCar,
		CarList:         carList,
		GraphData:       make([][]float64, 3, 100),
		GraphDataPoints: 100,
		Settings:        settings,
	}
	in := make(chan keys.Key)
	ch := make(chan fmtel.ForzaPacket)
	if err != nil {
		log.Error(err)
	}

	conn, err := net.ListenPacket("udp4", settings.UdpAddress)
	if err != nil {
		log.Error(err)
	}

	defer conn.Close()
	log.Debug("Starting server!", "address", settings.UdpAddress)

	go server.ReadPackets(conn, ch)
	go input.ListenForInput(in)
	if enableJson {
		go serveHTTP()
	}
	out.ClearScreen()
	var packet fmtel.ForzaPacket
	for {
		select {
		case key := <-in:
			{
				switch key.Code {
				case keys.CtrlT:
					{
						t := func() units.Temperature {
							if app.Settings.Temperature == units.CELSIUS {
								return units.FAHRENHEIT
							} else {
								return units.CELSIUS
							}
						}()
						app.Settings.Temperature = t
					}
				case keys.CtrlC, keys.Escape:
					{
						out.ExitAltScreen()
						restoreConsole()
						os.Exit(0)
					}
				default:
					{
						continue
					}
				}
			}
		case packet = <-ch:
			{
			}
			if !packet.IsPaused() {
				continue
			}

			if Pack.TimestampMS == packet.TimestampMS {
				continue
			}

			if cars.HasCarChanged(Pack.CarOrdinal, packet.CarOrdinal) {
				result := cars.SetCurrentCar(app.CarList, packet.CarOrdinal)
				app.CurrentCar = result
			}

			Pack = packet
			// area.Clear()
			layout := tui.Render(&packet, &app)
			if err != nil {
				log.Error(err)
			}

			// out.ClearScreen()
			out.MoveCursor(0, 0)
			out.WriteString(layout)
			// area.Update(layout)

		}
	}
}
