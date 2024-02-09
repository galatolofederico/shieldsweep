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

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/pkg/errors"
)

func get(httpc http.Client, path string) []byte {
	response, err := httpc.Get(path)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != http.StatusOK {
		panic(errors.Errorf("Error: %s\n", response.Status))
	}
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(errors.Wrapf(err, "Error reading response body from %s\n", path))
	}
	return resBody
}

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	sock := filepath.Join(home, "shsw.sock")
	_, err := os.Stat(sock)
	if err != nil {
		panic(errors.Wrapf(err, "Error checking for socket file: %v\n", sock))
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock)
			},
		},
	}

	command := flag.Args()

	switch command[0] {
	case "run":
		raw := get(httpc, "http://unix/run")
		var response messages.RunReply
		json.Unmarshal(raw, &response)
		if response.Started {
			color.Green("[+] Scan started")
		} else {
			color.Yellow("[-] Scan already running")
		}
	case "status":
		raw := get(httpc, "http://unix/status")
		var response messages.StatusReply
		json.Unmarshal(raw, &response)
		if response.Running {
			fdate := utils.ParseDate(response.StartedAt)
			color.White("[-] Scan running since " + fdate)
		} else {
			color.White("[-] SHSW is ready to scan")
		}
		for _, tool := range response.Tools {
			lastRun := utils.DaysAgo(tool.LastRun)
			lastLogChange := utils.DaysAgo(tool.LastLogChange)
			toolInfo := fmt.Sprintf("(last run: %s, last log change: %s)", lastRun, lastLogChange)
			switch tool.State {
			case "ready":
				color.Green("[-] " + tool.Name + " ready " + toolInfo)
			case "running":
				color.Green("[+] " + tool.Name + " running " + toolInfo)
			case "failed":
				color.Red("[-] " + tool.Name + " failed " + toolInfo)
			case "queued":
				color.Yellow("[+] " + tool.Name + " queued " + toolInfo)
			case "finished":
				color.Cyan("[+] " + tool.Name + " finished " + toolInfo)
			}
		}
	default:
		fmt.Println("Usage: shsw [run|status]")
	}
}
