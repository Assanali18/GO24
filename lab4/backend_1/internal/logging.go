package logging

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var Logger *logrus.Logger

func InitLogger() {

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			log.Fatalf("Не удалось создать папку logs: %v", err)
		}
	}

	logFile := &lumberjack.Logger{
		Filename:   "logs/application.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	Logger.SetOutput(multiWriter)

	Logger.SetLevel(logrus.InfoLevel)
}
