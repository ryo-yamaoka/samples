package main

import (
	"github.com/ryo-yamaoka/samples/slog-wrap-row-num/logger"
	"github.com/ryo-yamaoka/samples/slog-wrap-row-num/stdlog"
)

func main() {
	l := stdlog.New()
	logger.SetLogger(l)
	run()
}

func run() {
	logger.Info("info log")
	logger.Error("error log")
}
