package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/galatolofederico/shieldsweep/shsw/internal/engine"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	if !utils.IsRoot() {
		panic("shsw daemon must be run as root")
	}

	engine := engine.NewEngine(home)
	app := fiber.New(fiber.Config{
		ServerHeader: "shsw-daemon",
		AppName:      "shsw-daemon",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/run", func(c *fiber.Ctx) error {
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

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(messages.StatusReply{
			Running:   engine.IsRunning(),
			StartedAt: engine.GetStartedAt(),
			Tools:     engine.GetToolStates(),
		})
	})

	app.Get("/logs/:tool", func(c *fiber.Ctx) error {
		toolname := c.Params("tool")
		tool := engine.GetTool(toolname)
		if tool == nil {
			return c.Status(404).SendString("Tool not found")
		}
		return c.JSON(messages.LogsReply{
			Tool:            tool.Name,
			LatestLogChange: tool.State.LatestLogChange,
			LatestRun:       tool.State.LatestRun,
			Logs:            tool.GetLogs(),
		})
	})

	app.Get("/log/:tool/:id?", func(c *fiber.Ctx) error {
		toolname := c.Params("tool")
		logid, _ := strconv.Atoi(c.Params("id"))
		tool := engine.GetTool(toolname)
		if tool == nil {
			return c.Status(404).SendString("Tool not found")
		}
		logs := tool.GetLogs()
		if len(logs) <= 0 {
			return c.JSON(messages.LogReply{
				Tool:            tool.Name,
				State:           tool.State.State,
				LatestLogChange: "never",
				LatestRun:       tool.State.LatestRun,
				LatestError:     tool.GetLatestError(),
				Log:             "No logs available",
			})
		}
		if logid < 0 || logid >= len(logs) {
			logid = 0
		}
		log, err := os.ReadFile(logs[logid])
		if err != nil {
			log = []byte("Error reading log")
		}

		return c.JSON(messages.LogReply{
			Tool:            tool.Name,
			State:           tool.State.State,
			LatestLogChange: tool.State.LatestLogChange,
			LatestRun:       tool.State.LatestRun,
			LatestError:     tool.GetLatestError(),
			Log:             string(log),
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
	err = os.Chmod(sock, 0766)
	if err != nil {
		panic(err)
	}
	app.Listener(ln)
}
