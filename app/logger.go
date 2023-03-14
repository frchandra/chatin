package app

import (
	"github.com/frchandra/chatin/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

func NewLogger(config *config.AppConfig) *logrus.Logger {
	file, _ := os.OpenFile("./storage/logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	output := io.MultiWriter(os.Stdout, file)
	logger := logrus.New()
	logger.SetOutput(output)
	if config.IsProduction == "false" {
		logger.SetReportCaller(true)
		logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				//return frame.Function, fileName
				return frame.Function, fileName
			},
		})
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
