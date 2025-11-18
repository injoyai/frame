package frame

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "\x1b[36m", log.Ltime)
}
