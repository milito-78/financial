package main

import (
	"financial/bootstrap"
	"financial/config"
	"financial/io/telegram"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfg *config.App
)

func init() {
	var err error
	cfg, err = config.LoadConfig[config.App]()
	if err != nil {
		log.Fatalf("Error LoadConfig: %v", err.Error())
	}

	logLevel, _ := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(logLevel)
	if cfg.BotToken == "" {
		log.Fatal("Bot token cannot be empty!")
	}
}

func main() {
	bootstrap.InitDatabase()
	bootstrap.InitCache()
	bootstrap.InitDependencies()

	bot := telegram.NewTelegramBot(cfg.BotToken)

	go bot.StartBot()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	stop := <-signalForExit
	log.WithField("signal", stop).Info("GracefulStop signal Received ")
	log.Info("Waiting for all jobs to stop")
	bot.StopBot()
	log.Info("All jobs stop successfully")
}
