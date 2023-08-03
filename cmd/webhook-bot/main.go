package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/testit-tms/webhook-bot/internal/config"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/chat"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/company"
	"github.com/testit-tms/webhook-bot/internal/storage/postgres/owner"
	"github.com/testit-tms/webhook-bot/internal/transport/rest/send"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram"
	"github.com/testit-tms/webhook-bot/internal/transport/telegram/commands"
	"github.com/testit-tms/webhook-bot/internal/usecases"
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
	chatStorage := chat.New(db)

	regUsecases := registration.New(ownerStorage, companyStorage)
	registrator := commands.NewRegistrator(logger, regUsecases)

	companyUsesaces := usecases.NewCompanyUsecases(companyStorage, chatStorage)
	companyCommands := commands.NewCompanyCommands(companyUsesaces)

	chatUsesaces := usecases.NewChatUsecases(chatStorage, companyStorage)
	chatCommands := commands.NewChatCommands(chatUsesaces, companyUsesaces)

	bot, err := telegram.New(logger, cfg.TelegramBot.Token, registrator, companyCommands, chatCommands)
	if err != nil {
		logger.Error("cannot create telegram bot", err)
	}

	sendUsecases := usecases.NewSendMessageUsecases(logger, chatStorage, bot)
	handler := send.New(logger, sendUsecases)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// TODO: move to separate package
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/telegram", func(r chi.Router) {
		r.Post("/", handler)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()

	logger.Info("server is running")

	go func() {
		go bot.Run()
	}()

	logger.Info("telegram bot is running")

	<-done

	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", sl.Err(err))
		return
	}

	logger.Info("server stopped")
}
