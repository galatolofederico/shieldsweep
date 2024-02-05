package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"

	"github.com/galatolofederico/shieldsweep/shsw/internal/engine"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/gofiber/fiber/v3"
)

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	engine := engine.NewEngine(home)
	app := fiber.New()

	app.Get("/run", func(c fiber.Ctx) error {
		if engine.IsRunning() {
			return c.JSON(messages.RunReply{
				Started: false,
				Message: "Already running",
			})
		}
		go engine.Run()
		return c.JSON(messages.RunReply{
			Started: true,
			Message: "Scan started",
		})
	})

	app.Get("/status", func(c fiber.Ctx) error {
		return c.JSON(messages.StatusReply{
			Running:   engine.IsRunning(),
			StartedAt: engine.GetStartedAt(),
			Tools:     engine.GetToolStates(),
		})
	})

	sock := filepath.Join(home, "shsw.sock")
	if _, err := os.Stat(sock); err == nil {
		err = os.Remove(sock)
		if err != nil {
			panic(err)
		}
	}
	ln, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}
	app.Listener(ln)
}
