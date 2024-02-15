package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/galatolofederico/shieldsweep/shsw/internal/messages"
	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
)

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

	httpc := utils.GetUnixClient(sock)
	command := flag.Args()

	if len(command) < 1 {
		fmt.Println("Usage: shsw [run|status|list|log]")
		os.Exit(1)
	}

	switch command[0] {
	case "run":
		raw := utils.Get(httpc, "http://unix/run")
		var response messages.RunReply
		json.Unmarshal(raw, &response)
		if response.Started {
			color.Green("[+] Scan started")
		} else {
			color.Yellow("[-] Scan already running")
		}
	case "status":
		raw := utils.Get(httpc, "http://unix/status")
		var response messages.StatusReply
		json.Unmarshal(raw, &response)
		if response.Running {
			fdate := utils.ParseDate(response.StartedAt)
			color.White("[-] Scan running since " + fdate)
		} else {
			color.White("[-] SHSW is ready to scan")
		}
		for _, tool := range response.Tools {
			lastRun := utils.DaysAgo(tool.LatestRun)
			lastLogChange := utils.DaysAgo(tool.LatestLogChange)
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
	case "list":
		if len(command) < 2 {
			fmt.Println("Usage: shsw list <tool>")
			os.Exit(1)
		}
		tool := command[1]
		raw := utils.Get(httpc, "http://unix/logs/"+tool)
		var response messages.LogsReply
		json.Unmarshal(raw, &response)
		lastLogChange := utils.DaysAgo(response.LatestLogChange)
		color.Green("[!] Log for " + response.Tool)
		color.Green("[!] Latest log change: " + lastLogChange)
		color.White("[!] Available logs:")
		for i, log := range response.Logs {
			color.White(fmt.Sprintf("(id: %d) %v", i, log))
		}

	case "log":
		if len(command) < 3 {
			fmt.Println("Usage: shsw log <tool> <log-id>")
			fmt.Println("Use shsw list <tool> to get the log ids")
			os.Exit(1)
		}
		tool := command[1]
		logid := command[2]
		raw := utils.Get(httpc, "http://unix/log/"+tool+"/"+logid)
		var response messages.LogReply
		json.Unmarshal(raw, &response)
		lastLogChange := utils.DaysAgo(response.LatestLogChange)
		color.Green("[!] Log for " + response.Tool)
		color.Green("[!] Latest log change: " + lastLogChange)
		color.White(response.Log)
		if response.LatestError != "" {
			color.Red("[!] Error log found")
			color.Red(response.LatestError)
		}
	default:
		fmt.Println("Usage: shsw [run|status|list|log]")
	}
}
