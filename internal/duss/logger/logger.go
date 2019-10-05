// Package logger provides initializer and configuration for global Logging
package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// InitLogger sets the default configuration of logger with timestamp
// and redirects every output to stdout
func InitLogger() {
	defer log.Info("Logger configured successfully")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout) // everything to stdout
}
