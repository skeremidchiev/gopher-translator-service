package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/skeremidchiev/gopher-translator-service/app/configurate"
)

func getConfiguration() configurate.Config {
	configFilePath := os.Getenv("TRANSLATION_CONFIG")

	service, err := configurate.NewConfig(configFilePath)
	if err != nil {
		log.Fatal("[Main] ", err)
	}

	log.Info("[Main] Staking configuration acquired")

	return service
}
