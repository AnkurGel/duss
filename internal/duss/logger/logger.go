package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	defer log.Info("Logger configured successfully")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout) // everything to stdout
}
