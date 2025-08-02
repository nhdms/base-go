package logger

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
	//lg "go-micro.dev/v5/logger"
)

var DefaultLogger *ATLogger = defaultLogger()

func defaultLogger() *ATLogger {
	l, e := zap.NewDevelopment()
	if e != nil {
		log.Fatal("Failed to initialize logger:", e)
	}

	return &ATLogger{
		SugaredLogger: l.Sugar(),
	}
}

type ATLogger struct {
	*zap.SugaredLogger
}

type ctxKey struct{}

var logger *ATLogger

func InitLogger() {
	configLevel := viper.GetString("logger.level")
	log.Println("Initializing logger with level:", configLevel)

	lv, err := zap.ParseAtomicLevel(configLevel)
	if err != nil {
		lv = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	//sync.OnceFunc(func() {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = lv
	defaultConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	lg, err := defaultConfig.Build()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}

	DefaultLogger = &ATLogger{
		SugaredLogger: lg.Sugar(),
	}

	_ = lg.Sync() // flushes buffer, if any
}

func FromCtx(ctx context.Context) *ATLogger {
	if l, ok := ctx.Value(ctxKey{}).(*ATLogger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return &ATLogger{SugaredLogger: zap.NewNop().Sugar()}
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
