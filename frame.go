package frame

import (
	"fmt"
	"time"
)

const (
	DefaultPort = 8080
)

type Logger interface {
	Printf(format string, v ...any) (int, error)
	Println(v ...any) (int, error)
}

var _ Logger = (*log)(nil)

func NewLogger() Logger {
	return &log{}
}

type log struct{}

func (this *log) Printf(format string, v ...any) (int, error) {
	return fmt.Printf("\x1b[36m"+time.Now().Format(time.TimeOnly)+" "+format+"\x1b[0m", v...)
}

func (this *log) Println(v ...any) (int, error) {
	s := fmt.Sprint(v...)
	return fmt.Println("\x1b[36m", time.Now().Format(time.TimeOnly), " ", s, "\x1b[0m")
}
