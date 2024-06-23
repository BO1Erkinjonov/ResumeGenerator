package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"resume-generator/internal/app"
	configpkg "resume-generator/internal/pkg/config"
	"syscall"
)

func main() {
	// config
	config, err := configpkg.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// app
	app, err := app.NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	// run application
	go func() {
		app.Logger.Info("Listen: ", zap.String("address", "localhost"+config.Server.Port))
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	app.Logger.Info("api gateway service stops")
	app.Stop()
}
