package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sudhanshu042004/sandbox/internal/config"
	"github.com/sudhanshu042004/sandbox/internal/http/handlers/auth"
	"github.com/sudhanshu042004/sandbox/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//db setup
	storage, err := sqlite.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage intialised", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/login", auth.Login(storage))
	router.HandleFunc("POST /api/signup", auth.SignUp(storage))
	//setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()
	<-done

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
