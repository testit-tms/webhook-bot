package main

import (
	"log"
	"os"

	"github.com/testit-tms/webhook-bot/internal/config"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/chat"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/company"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/owner"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram/commands"
	"github.com/testit-tms/webhook-bot/internal/usecases/registration"
	"github.com/testit-tms/webhook-bot/pkg/database"
	"github.com/testit-tms/webhook-bot/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Initialize(cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		logger.Error("cannot initialize database", sl.Err(err))
		os.Exit(1)
	}

	ownerStorage := owner.New(db)
	companyStorage := company.New(db)
	_ = chat.New(db)

	regUsecases := registration.New(ownerStorage, companyStorage)
	registrator := commands.NewRegistrator(logger, regUsecases)
	bot, err := telegram.New(logger, cfg.TelegramBot.Token, registrator)
	if err != nil {
		logger.Error("cannot create telegram bot", err)
	}

	bot.Run()
}
