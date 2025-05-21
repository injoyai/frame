package frame

import (
	"github.com/fatih/color"
	"github.com/injoyai/logs"
)

const (
	DefaultPort = 8080
)

var (
	Log Logger = logs.NewEntity("").SetSelfLevel(logs.LevelInfo).SetColor(color.FgCyan).SetFormatter(logs.TimeFormatter)
)

type Logger interface {
	Printf(format string, v ...any) (int, error)
	Println(v ...any) (int, error)
}
