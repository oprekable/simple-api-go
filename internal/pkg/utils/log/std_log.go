package log

import (
	stdlog "log"

	"github.com/rs/zerolog"
)

type StdLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

var _ StdLogger = &stdlog.Logger{}

type stdLogger struct {
	log zerolog.Logger
}

func (s *stdLogger) Fatal(_ ...interface{}) {
	panic("implement me")
}

func (s *stdLogger) Fatalf(_ string, _ ...interface{}) {
	panic("implement me")
}

func (s *stdLogger) Print(_ ...interface{}) {
	panic("implement me")
}

func (s *stdLogger) Println(_ ...interface{}) {
	panic("implement me")
}

func (s *stdLogger) Printf(_ string, _ ...interface{}) {
	panic("implement me")
}

func StandardLogger(log zerolog.Logger) StdLogger {
	return &stdLogger{log}
}
