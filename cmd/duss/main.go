package main

import (
	"bytes"
	"fmt"
	"github.com/ankurgel/duss/internal/duss/logger"
	"github.com/ankurgel/duss/internal/duss/server"
	"github.com/ankurgel/duss/internal/duss/store"
	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
)

func main() {
	log.Println("Starting DUSS")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	logger.InitLogger()
	readConfigs()

	s := store.InitStore()
	defer s.Close()

	h := server.InitServer(s)
	h.SetHandlers()
	go func() {
		if err := h.Listen(viper.GetString("ListenAddr")); err != nil {
			// TODO: handle actual failure with panic
			log.Error(err)
		}
	}()

	<- quit
	s.Close()
	h.Close()
	log.Println("Shutting down")
}


func readConfigs() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var configBox = packr.New("Configs", "../../configs")
	var configFilePath = viper.GetString("DUSS_CONFIG_PATH")

	var yamlContent []byte
	var err error
	if configFilePath == "" {
		configFilePath = "config.yaml"
		yamlContent, err = configBox.Find(configFilePath)
	} else {
		yamlContent, err = ioutil.ReadFile(configFilePath)
	}

	if err != nil {
		log.Error(err)
		panic(fmt.Errorf("error in Parsing Configration(): %s", err))
	}

	if err = viper.ReadConfig(bytes.NewBuffer(yamlContent)); err != nil {
		log.Error(err)
		panic(fmt.Errorf("error in ReadConfigs(): %s", err))
	}
	log.Info(viper.GetString("Environment"), " configuration set successfully")
}
