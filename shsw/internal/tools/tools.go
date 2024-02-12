package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/galatolofederico/shieldsweep/shsw/internal/utils"
	"github.com/pkg/errors"
)

const (
	Ready    = "ready"
	Running  = "running"
	Queued   = "queued"
	Finished = "finished"
	Failed   = "failed"
)

type ToolState struct {
	LastRun       string
	LastLogChange string
	LastLogHash   string
	LastError     string
	State         string
}

type ToolRunner interface {
	Check() bool
	Run(config Tool) error
}

type ToolConfig struct {
	Name    string
	Enabled bool
	Config  json.RawMessage
}

type Tool struct {
	State          ToolState
	Runner         ToolRunner
	Name           string
	LogsPath       string
	CurrentLogFile string
	StateFile      string
}

type ToolResult struct {
	Name     string
	IsLogNew bool
	Error    error
}

func (tool *Tool) Run(ch chan<- ToolResult) {
	now := time.Now()
	result := ToolResult{Name: tool.Name, IsLogNew: false, Error: nil}
	tool.CurrentLogFile = filepath.Join(tool.LogsPath, now.Format("2006-01-02 15:04:05")+".txt")

	utils.CheckPathForFile(tool.CurrentLogFile)

	tool.State.LastRun = now.Format(time.RFC3339)
	tool.State.State = Running
	err := tool.Runner.Run(*tool)

	newLogHash := "none"
	if err == nil {
		exists := utils.FileExists(tool.CurrentLogFile)
		if exists {
			newLogHash, err = utils.SHA256File(tool.CurrentLogFile)
		} else {
			err = errors.Errorf("Log file not found: %v", tool.CurrentLogFile)
		}
	}

	if err != nil {
		tool.State.LastError = err.Error()
		tool.State.State = Failed
		result.Error = err
	} else {
		if newLogHash != tool.State.LastLogHash {
			tool.State.LastLogChange = time.Now().Format(time.RFC3339)
			tool.State.LastLogHash = newLogHash
			result.IsLogNew = true
		}
		tool.State.LastError = ""
		tool.State.State = Finished
	}

	tool.Save()
	ch <- result
}

func (tool *Tool) Load() {
	utils.CheckPathForFile(tool.StateFile)
	if _, err := os.Stat(tool.StateFile); os.IsNotExist(err) {
		tool.State = ToolState{
			LastRun:       "never",
			LastLogChange: "never",
			LastLogHash:   "none",
			LastError:     "",
			State:         Ready,
		}
		tool.Save()
	}
	dat, err := os.ReadFile(tool.StateFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &tool.State)
	if err != nil {
		panic(err)
	}
}

// TODO: se non riesce a salvare il file di stato deve davvero panicare?
func (tool *Tool) Save() {
	utils.CheckPathForFile(tool.StateFile)
	encoded, err := json.Marshal(tool.State)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(tool.StateFile, encoded, 0644)
	if err != nil {
		panic(err)
	}
}

func (tool *Tool) GetLog() string {
	utils.CheckPathForFile(tool.CurrentLogFile)
	dat, err := os.ReadFile(tool.CurrentLogFile)
	if err != nil {
		return "Log file not found"
	}
	return string(dat)
}

func (tool *Tool) GetLatestLog() string {
	files, err := os.ReadDir(tool.LogsPath)
	if err != nil {
		return "Error reading log files"
	}

	var dates []time.Time
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			dateStr := strings.TrimSuffix(filename, filepath.Ext(filename))
			date, err := time.Parse("2006-01-02 15:04:05", dateStr)
			if err == nil {
				dates = append(dates, date)
			}
		}
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})

	if len(dates) > 0 {
		latestDate := dates[0]
		latestDateStr := latestDate.Format("2006-01-02 15:04:05")
		return filepath.Join(tool.LogsPath, latestDateStr+".log")
	}

	return ""
}

func (tool *Tool) GetLastError() string {
	return tool.State.LastError
}
