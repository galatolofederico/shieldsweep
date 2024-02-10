package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/pkg/errors"
)

type ToolStatus struct {
	Name          string
	State         string
	LastRun       string
	LastLogChange string
}

type StatusPageData struct {
	Running     bool
	StartedAt   string
	ToolsStatus []ToolStatus
}

var sock string

func get(httpc http.Client, path string) []byte {
	response, err := httpc.Get(path)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(response.Body)
		panic(fmt.Errorf("Error: %s\n%s", response.Status, resBody))
	}
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading response body from %s\n", path))
	}
	return resBody
}

func getHTTPClient() http.Client {
	return http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock)
			},
		},
	}
}

func statusHandler(c *fiber.Ctx) error {
	httpc := getHTTPClient()

	raw := get(httpc, "http://unix/status")
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
			Name:          tool.Name,
			State:         tool.State,
			LastRun:       tool.LastRun,
			LastLogChange: tool.LastLogChange,
		}
	}

	return c.Render("status", fiber.Map{
		"Running":     statusData.Running,
		"StartedAt":   statusData.StartedAt,
		"ToolsStatus": statusData.ToolsStatus,
	}, "status")
}

func startScanHandler(c *fiber.Ctx) error {
	httpc := getHTTPClient()

	raw := get(httpc, "http://unix/run")
	var response messages.RunReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing run response")
	}

	return c.Redirect("/")

}

func toolDetailHandler(c *fiber.Ctx) error {
	toolName := c.Params("toolName")
	httpc := getHTTPClient()

	raw := get(httpc, "http://unix/log/"+toolName)
	var response messages.LogReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing log response")
	}

	return c.Render("tool_detail", fiber.Map{
		"Name":          toolName,
		"LastRun":       response.LastRun,
		"LastLogChange": response.LastLogChange,
		"Logs":          response.Log,
		"Errors":        response.LastError,
	}, "tool_detail")
}

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	sock = filepath.Join(home, "shsw.sock")
	_, err := os.Stat(sock)
	if err != nil {
		panic(errors.Wrapf(err, "Error checking for socket file: %v\n", sock))
	}

	engine := html.New("./views", ".html")

	engine.AddFunc("stateToClass", func(state string) string {
		switch state {
		case "ready":
			return "table-success"
		case "running":
			return "table-primary"
		case "queued":
			return "table-warning"
		case "failed":
			return "table-danger"
		default:
			return ""
		}
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", statusHandler)
	app.Post("/start-scan", startScanHandler)
	app.Get("/tool/:toolName", toolDetailHandler)

	app.Listen(":3000")
}
