package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/muesli/termenv"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard/keys"
	flag "github.com/spf13/pflag"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cmd/fmtui/input"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/charmbracelet/log"
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/server"
	"github.com/stelmanjones/fmtel/units"
)

var Pack = fmtel.DefaultForzaPacket

type App struct {
	Settings   Settings
	CarList    []cars.Car
	CurrentCar cars.Car
}

type Settings struct {
	Temperature units.Temperature
	UdpAddress  string
}

// HACK: Move these to the settings struct?
var (
	temp       string
	udpAddress string
	enableJson bool
	enableSSE  bool
	baseUrl    string
	noUi       bool
)

// TODO: Rename this function.
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

func serveHTTP(address string) {
	if enableSSE {
		s := sse.NewServer(&sse.Options{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		})
		defer s.Shutdown()
		http.Handle("/sse", s)

		// HACK: There is probably a better way to do this loop.
		go func() {
			ticker := time.Tick(200 * time.Millisecond)
			for {
				data, err := Pack.ToJson()
				if err != nil {
					log.Error(err)
					return
				}

				<-ticker
				if s.ClientCount() == 0 {
					log.Debug("No clients connected")
					time.Sleep(1 * time.Second)
					continue
				}
				s.SendMessage("/sse", sse.SimpleMessage(string(data)))
			}
		}()
	}
	if enableJson {
		http.HandleFunc("/json", responder)
	}

	log.Debugf("Telemetry Server started at %s", baseUrl)
	log.Fatal(http.ListenAndServe(address, nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	flag.StringVar(&temp, "temp", "celsius", "Set temperature unit.")
	flag.StringVar(&udpAddress, "udp-addr", ":7777", "Set UDP connection address.")
	flag.StringVar(&baseUrl, "base-url", ":9999", "Set telemetry server address.")
	flag.BoolVar(&enableJson, "json", false, "Enable JSON endpoint.")
	flag.BoolVar(&enableSSE, "sse", false, "Enable SSE endpoint.")
	flag.BoolVar(&noUi, "no-ui", false, "Run without TUI.")
	flag.Lookup("json").NoOptDefVal = "true"
	flag.Lookup("sse").NoOptDefVal = "true"
	flag.Lookup("no-ui").NoOptDefVal = "true"
	flag.Parse()

	out := termenv.DefaultOutput()

	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(out)
	if err != nil {
		panic(err)
	}
	if !noUi {
		cursor.Hide()
		out.AltScreen()

		defer cursor.Show()
	} else {
		log.SetLevel(log.DebugLevel)
	}

	settings := types.Settings{
		Temperature: units.TempFromString(temp),
		UdpAddress:  udpAddress,
	}
	carList, err := cars.ReadCarList("cars.json")
	if err != nil {
		log.Error(err)
	}
	app := types.App{
		CurrentCar: cars.DefaultCar,
		CarList:    carList,
		Settings:   settings,
	}
	in := make(chan keys.Key)
	ch := make(chan fmtel.ForzaPacket)
	if err != nil {
		log.Error(err)
	}

	shutdown := func() {
		out.ExitAltScreen()
		restoreConsole()
		close(in)
		close(ch)
		os.Exit(0)
	}

	conn, err := net.ListenPacket("udp4", settings.UdpAddress)
	if err != nil {
		log.Error(err)
	}

	defer conn.Close()
	log.Debug("Starting server!", "address", settings.UdpAddress)

	go server.ReadPackets(conn, ch)
	go input.ListenForInput(in)
	if enableJson || enableSSE {
		go serveHTTP(baseUrl)
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
						shutdown()
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

			Pack = packet

			if !noUi {
				if cars.HasCarChanged(Pack.CarOrdinal, packet.CarOrdinal) {
					result := cars.SetCurrentCar(app.CarList, packet.CarOrdinal)
					app.CurrentCar = result
				}

				layout := tui.Render(&packet, &app)
				if err != nil {
					log.Error(err)
				}

				out.MoveCursor(0, 0)
				out.WriteString(layout)
			}

		}
	}
}
