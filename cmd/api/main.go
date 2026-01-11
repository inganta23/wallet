package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inganta23/wallet/internal/config"
	"github.com/inganta23/wallet/internal/handler"
	"github.com/inganta23/wallet/internal/repository"
	"github.com/inganta23/wallet/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cfg.DBMaxConn)
	db.SetMaxIdleConns(cfg.DBMaxConn)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		slog.Error("DB unreachable", "error", err)
		os.Exit(1)
	}

	repo := repository.NewPostgresRepo(db)
	svc := service.NewWalletService(repo)
	h := &handler.WalletHandler{Service: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /balance", h.GetBalance)
	mux.HandleFunc("POST /withdraw", h.Withdraw)

	server := &http.Server{Addr: cfg.ServerPort, Handler: mux}

	go func() {
		slog.Info("Server starting", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	slog.Info("Server stopped")
}