package main

import (
	"fmt"
	"github.com/ankurgel/duss/internal/duss/logger"
	"github.com/ankurgel/duss/internal/duss/server"
	"github.com/ankurgel/duss/internal/duss/store"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func main() {
	log.Println("Starting DUSS")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	logger.InitLogger()
	ReadConfigs()

	//db := store.InitDB()
	//defer db.Close()
	s := store.InitStore()
	defer s.Close()

	h := server.InitServer()
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


func ReadConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
		panic(fmt.Errorf("Error in ReadConfigs(): %s\n", err))
	}
	log.Info("Configuration set successfully")
}