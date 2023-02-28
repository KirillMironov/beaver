package logger

import (
	"os"
	"path/filepath"
	"time"

	"golang.org/x/exp/slog"
)

func New(level slog.Level) *slog.Logger {
	options := slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.TimeValue(time.Now().UTC())
			case slog.SourceKey:
				a.Value = slog.StringValue(filepath.Base(a.Value.String()))
			}

			return a
		},
	}

	return slog.New(options.NewJSONHandler(os.Stdout))
}
