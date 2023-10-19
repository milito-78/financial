package main

import (
	"context"
	"financial/application"
	"financial/bootstrap"
	"financial/config"
	"financial/infrastructure/cache"
	"financial/io/telegram"
	"github.com/golobby/container/v3"
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

	ctx := context.TODO()
	var c cache.ICache
	_ = container.Resolve(&c)

	var userService application.IUserService
	_ = container.Resolve(&userService)

	bot := telegram.NewTelegramBot(ctx, c, cfg.BotToken, userService)

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
