package newbot

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type NoOpLogger struct{}

func (NoOpLogger) Printf(format string, v ...any) {}
func (NoOpLogger) Println(v ...any)               {}
