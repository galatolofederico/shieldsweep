package main

import (
	"embed"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
)

type ToolStatus struct {
	Name            string
	State           string
	LatestRun       string
	LatestLogChange string
}

type StatusPageData struct {
	Running     bool
	StartedAt   string
	ToolsStatus []ToolStatus
}

var httpc http.Client

func statusHandler(c *fiber.Ctx) error {
	raw := utils.Get(httpc, "http://unix/status")
	var response messages.StatusReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing status response")
	}

	statusData := StatusPageData{
		Running:     response.Running,
		StartedAt:   response.StartedAt,
		ToolsStatus: make([]ToolStatus, len(response.Tools)),
	}

	for i, tool := range response.Tools {
		statusData.ToolsStatus[i] = ToolStatus{
			Name:            tool.Name,
			State:           tool.State,
			LatestRun:       utils.DaysAgo(tool.LatestRun),
			LatestLogChange: utils.DaysAgo(tool.LatestLogChange),
		}
	}

	return c.Render("views/status", fiber.Map{
		"Running":     statusData.Running,
		"StartedAt":   statusData.StartedAt,
		"ToolsStatus": statusData.ToolsStatus,
	}, "views/status")
}

func startScanHandler(c *fiber.Ctx) error {
	raw := utils.Get(httpc, "http://unix/run")
	var response messages.RunReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing run response")
	}

	return c.Redirect("/")

}

func toolDetailHandler(c *fiber.Ctx) error {
	toolName := c.Params("toolName")
	logId, _ := strconv.Atoi(c.Params("logId"))
	raw := utils.Get(httpc, "http://unix/log/"+toolName+"/"+strconv.Itoa(logId))
	var logResponse messages.LogReply
	if err := json.Unmarshal(raw, &logResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing log response")
	}
	raw = utils.Get(httpc, "http://unix/logs/"+toolName)
	var logsResponse messages.LogsReply
	if err := json.Unmarshal(raw, &logsResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing logs response")
	}

	return c.Render("views/tool_detail", fiber.Map{
		"Name":            toolName,
		"State":           logResponse.State,
		"LatestRun":       utils.DaysAgo(logResponse.LatestRun),
		"LatestLogChange": utils.DaysAgo(logResponse.LatestLogChange),
		"LogDate":         logResponse.LogDate,
		"Logs":            logsResponse.Logs,
		"Log":             logResponse.Log,
		"Error":           logResponse.LatestError,
	}, "views/tool_detail")
}

//go:embed views/*
var embedDirViews embed.FS

//go:embed assets/*
var embedDirAssets embed.FS

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	sock := filepath.Join(home, "shsw.sock")
	_, err := os.Stat(sock)
	if err != nil {
		color.Red("[!] Error: %v\n", err)
		color.Red("[!] Is the daemon running?\n")
	}
	httpc = utils.GetUnixClient(sock)
	utils.Get(httpc, "http://unix/health")

	engine := html.NewFileSystem(http.FS(embedDirViews), ".html")

	engine.AddFunc("stateToClass", func(state string) string {
		switch state {
		case "ready":
			return "primary"
		case "running":
			return "success"
		case "queued":
			return "warning"
		case "failed":
			return "danger"
		case "finished":
			return "info"
		default:
			return ""
		}
	})

	engine.AddFunc("fileToDate", func(fname string) string {
		parts := strings.Split(fname, "/")
		date := strings.Split(parts[len(parts)-1], ".")[0]
		return date
	})

	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "shsw-web",
		AppName:      "shsw-web",
	})

	app.Get("/", statusHandler)
	app.Post("/start-scan", startScanHandler)
	app.Get("/tool/:toolName/:logId", toolDetailHandler)

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirAssets),
		PathPrefix: "assets",
		Browse:     true,
	}))
	app.Listen(":3000")
}
