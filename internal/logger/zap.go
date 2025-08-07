package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Initialize(level zapcore.Level) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(zapcore.Lock(&zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(os.Stdout), // bisa diganti ke file
			Size:          4096,                       // buffer 4KB
			FlushInterval: 100 * time.Millisecond,
		})),
		level,
	)

	Log = zap.New(core, zap.AddCaller())
}

func Sync() {
	_ = Log.Sync()
}
