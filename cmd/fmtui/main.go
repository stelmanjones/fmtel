package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard/keys"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	flag "github.com/spf13/pflag"
	"github.com/stelmanjones/fmtel"
	"github.com/stelmanjones/fmtel/cars"
	"github.com/stelmanjones/fmtel/cmd/fmtui/input"
	"github.com/stelmanjones/fmtel/cmd/fmtui/pyro"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui"
	"github.com/stelmanjones/fmtel/cmd/fmtui/tui/views"
	"github.com/stelmanjones/fmtel/cmd/fmtui/types"
	"github.com/stelmanjones/fmtel/server"
	"github.com/stelmanjones/fmtel/units"
)

// Pack is just a global packet
var (
	Pack = fmtel.DefaultForzaPacket
	Dyno = tui.NewDynoView()
)

var (
	enableProfiling bool
	udpAddress      string
	enableJSON      bool
	enableSSE       bool
	baseURL         string
	UI              bool
	refreshRate     int
)

func jsonResponder(w http.ResponseWriter, r *http.Request) {
	data, err := Pack.ToJSON()
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
				data, err := Pack.ToJSON()
				if err != nil {
					log.Error(err)
					return
				}

				<-ticker
				if s.ClientCount() == 0 {
					log.Debug("No clients connected")
					time.Sleep(500 * time.Millisecond)
					continue
				}
				s.SendMessage("/sse", sse.SimpleMessage(string(data)))
			}
		}()
	}
	if enableJSON {
		http.HandleFunc("/json", jsonResponder)
	}

	log.Debugf("Telemetry Server started at %s", baseURL)
	log.Fatal(http.ListenAndServe(address, nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	flag.StringVar(&udpAddress, "udp-addr", ":7777", "Set UDP connection address.")
	flag.StringVar(&baseURL, "base-url", ":9999", "Set telemetry server address.")
	flag.BoolVar(&enableJSON, "json", false, "Enable JSON endpoint.")
	flag.BoolVar(&enableSSE, "sse", false, "Enable SSE endpoint.")
	flag.BoolVar(&UI, "tui", false, "Run with TUI.")
	flag.BoolVarP(&enableProfiling, "profiling", "p", false, "Enables pyroscope profiling for this app.")
	flag.IntVar(&refreshRate, "refresh", 60, "Set refresh rate per second.")
	flag.Lookup("json").NoOptDefVal = "true"
	flag.Lookup("sse").NoOptDefVal = "true"
	flag.Lookup("tui").NoOptDefVal = "true"
	flag.Lookup("profiling").NoOptDefVal = "true"

	flag.Parse()

	if enableProfiling {
		pyro.RunProfiling()
	}

	settings := types.Settings{
		Temperature:     units.CELSIUS,
		UdpAddress:      udpAddress,
		EndpointAddress: baseURL,
		EnableJSON:      enableJSON,
		EnableSSE:       enableSSE,
		RefreshRate:     refreshRate,
	}
	carList, err := cars.ReadCarList("cars.json")
	if err != nil {
		log.Error(err)
	}

	out := termenv.DefaultOutput()

	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(out)
	if err != nil {
		panic(err)
	}
	if UI {
		cursor.Hide()
		out.AltScreen()

		defer cursor.Show()
	} else {
		log.SetLevel(log.DebugLevel)
	}

	app := types.App{
		CurrentCar:  cars.DefaultCar,
		CarList:     carList,
		Settings:    settings,
		CurrentView: views.Home,
	}

	in := make(chan keys.Key)
	ch := make(chan fmtel.ForzaPacket)

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

	go server.ReadPackets(conn, ch, app.Settings.RefreshRate)
	go input.ListenForInput(in)
	if app.Settings.EnableJSON || app.Settings.EnableSSE {
		log.Debug("Enabled features", "SSE", app.Settings.EnableSSE, "JSON", app.Settings.EnableJSON)
		go serveHTTP(app.Settings.EndpointAddress)
	}
	out.ClearScreen()
	var packet fmtel.ForzaPacket
	for {
		select {
		case key := <-in:
			{
				switch key.Code {
				case keys.RuneKey:
					{
						switch key.String() {
						case "q":
							{
								shutdown()
							}
						case "t":
							{
								t := func() units.Temperature {
									if app.Settings.Temperature == units.CELSIUS {
										return units.FAHRENHEIT
									}
									return units.CELSIUS
								}()
								app.Settings.Temperature = t

							}
						case "r":
							{
								if app.CurrentView == views.Dyno {
									Dyno.Reset()
								}
							}
						case "d":
							{
								out.ClearScreen()
								if app.CurrentView != views.Dyno {
									app.CurrentView = views.Dyno
								} else {
									app.CurrentView = views.Home
								}
							}
						default:
							{
							}
						}
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

			if UI {
				if cars.HasCarChanged(Pack.CarOrdinal, packet.CarOrdinal) {
					result := cars.SetCurrentCar(app.CarList, packet.CarOrdinal)
					app.CurrentCar = result
				}
				Dyno.Update(&packet)

				// TODO: Split rendering to a separate function that switches on view.
				switch app.CurrentView {

				case views.Dyno:
					{

						view, err := tui.RenderDynoView(Dyno, &app.CurrentCar)
						if err != nil {
							log.Error(err)
						}
						out.MoveCursor(0, 0)
						out.WriteString(view)

					}
				case views.Home:
					{

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
	}
}
