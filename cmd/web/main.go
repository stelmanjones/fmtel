package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/spf13/pflag"
	"github.com/stelmanjones/fmtel/internal"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/storage/bbolt/v2"
)

var (
	address string
	store   *bbolt.Storage
)

func main() {
	store = bbolt.New()
	defer store.Close()

	pflag.StringVarP(&address, "address", "a", ":1337", "set server address.")
	pflag.Parse()

	log.Info("Starting! 🏁")

	engine := handlebars.New("cmd/web/views", ".hbs")

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("FMTEL %s", internal.VERSION),
		Views:   engine,
	})
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("ENV") == "development"
		},
	}))

	app.Get("/metrics", monitor.New(monitor.Config{Title: "FMTEL Web Metrics"}))

	app.Get("/", func(c *fiber.Ctx) error {
		k, err := store.Get("count")
		if err != nil {
			log.Error(err)
		}
		current, err := strconv.Atoi(string(k))
		if err != nil {
			log.Error(err)
		}
		current++
		store.Set("count", []byte(strconv.Itoa(current)), 0)
		return c.Render("main", fiber.Map{
			"count": current,
			"title": "FMTEL Dashboard",
		}, "layouts/shell")
	})
	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			var i int
			log.Info("SSE Client connected!")
			for {
				i++
				msg := fmt.Sprintf("%d", i)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				log.Debug("SSE", "count", i)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					log.Errorf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(1 * time.Second)
			}
		}))
		return nil
	})

	log.Fatal(app.Listen(address))
	log.Info("Bye!")
}
