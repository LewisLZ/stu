package web

import (
	log "github.com/sirupsen/logrus"
)

var (
	//binary version
	Version = "unknown"
)

func init() {
	log.Infof("Version: %v", Version)
}
