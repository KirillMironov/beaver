package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*zap.SugaredLogger)(nil)

type Logger interface {
	Info(...any)
	Infof(string, ...any)
	Error(...any)
	Errorf(string, ...any)
}

func New() (*zap.SugaredLogger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.ConsoleSeparator = " "

	config := zap.NewProductionConfig()
	config.EncoderConfig = encoderConfig
	config.Encoding = "console"
	config.DisableCaller = true

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("01-02|15:04:05.000") + "]")
}
