package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
