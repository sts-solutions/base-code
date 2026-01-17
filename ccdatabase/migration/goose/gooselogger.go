package goose

import (
	"context"
	"fmt"

	"github.com/sts-solutions/base-code/cclogger"
)

type gooseLogger struct {
	logger cclogger.Logger
	ctx    context.Context
}

func newGooseLogger(logger cclogger.Logger) gooseLogger {
	return gooseLogger{
		logger: logger,
		ctx:    context.Background(),
	}
}

func (g gooseLogger) Print(v ...interface{}) {
	g.logger.Info(g.ctx, fmt.Sprint(v...))
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	g.logger.Info(g.ctx, msg)
}

func (g gooseLogger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...) + "\n"
	g.logger.Info(g.ctx, msg)
}

func (g gooseLogger) Fatal(v ...interface{}) {
	g.logger.Fatal(g.ctx, fmt.Sprint(v...))
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	g.logger.Fatal(g.ctx, msg)
}

func (g gooseLogger) Fatalln(v ...interface{}) {
	msg := fmt.Sprint(v...) + "\n"
	g.logger.Fatal(g.ctx, msg)
}
