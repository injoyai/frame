package frame

const (
	DefaultPort = 8080
)

type Logger interface {
	Printf(format string, v ...any) (int, error)
	Println(v ...any) (int, error)
}
