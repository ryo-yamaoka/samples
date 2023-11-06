package stdlog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

type StdLogger struct {
	stdLogger *slog.Logger
	errLogger *slog.Logger
}

func New() *StdLogger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(_ []string, att slog.Attr) slog.Attr {
			if att.Key == slog.SourceKey {
				const skip = 8
				_, file, line, ok := runtime.Caller(skip)
				if !ok {
					return att
				}
				v := fmt.Sprintf("%s:%d", filepath.Base(file), line)
				att.Value = slog.StringValue(v)
			}
			return att
		},
	}
	sl := slog.New(slog.NewTextHandler(os.Stdout, opts))
	el := slog.New(slog.NewTextHandler(os.Stderr, opts))

	return &StdLogger{
		stdLogger: sl,
		errLogger: el,
	}
}

func (ml *StdLogger) Info(format string, v ...interface{}) {
	ml.stdLogger.Info(format, v...)
}

func (ml *StdLogger) Error(format string, v ...interface{}) {
	ml.errLogger.Error(format, v...)
}
