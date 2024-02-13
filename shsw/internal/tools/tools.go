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
	LatestRun       string
	LatestLogChange string
	LatestLogHash   string
	LatestError     string
	State           string
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

	tool.State.LatestRun = now.Format(time.RFC3339)
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
		tool.State.LatestError = err.Error()
		tool.State.State = Failed
		result.Error = err
	} else {
		if newLogHash != tool.State.LatestLogHash {
			tool.State.LatestLogChange = time.Now().Format(time.RFC3339)
			tool.State.LatestLogHash = newLogHash
			result.IsLogNew = true
		}
		tool.State.LatestError = ""
		tool.State.State = Finished
	}

	tool.Save()
	ch <- result
}

func (tool *Tool) Load() {
	utils.CheckPathForFile(tool.StateFile)
	if _, err := os.Stat(tool.StateFile); os.IsNotExist(err) {
		tool.State = ToolState{
			LatestRun:       "never",
			LatestLogChange: "never",
			LatestLogHash:   "none",
			LatestError:     "",
			State:           Ready,
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

func (tool *Tool) GetLogs() []string {
	files, err := os.ReadDir(tool.LogsPath)
	if err != nil {
		panic(errors.Wrap(err, "Cant read the logs path"))
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

	var logs []string
	for _, date := range dates {
		stringDate := date.Format("2006-01-02 15:04:05")
		logs = append(logs, filepath.Join(tool.LogsPath, stringDate+".txt"))
	}

	return logs
}

func (tool *Tool) GetLogById(index int) string {
	logs := tool.GetLogs()
	if index < len(logs) {
		return logs[index]
	}

	return ""
}

func (tool *Tool) GetLatestLog() string {
	return tool.GetLogById(0)
}

func (tool *Tool) GetLatestError() string {
	return tool.State.LatestError
}
