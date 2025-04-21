package main

import (
	"context"
	"go-svc-gophermart/internal/auth"
	"go-svc-gophermart/internal/config"
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/repositories"
	"go-svc-gophermart/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func StartServer(serverURL string, router *chi.Mux) error {
	var err error
	server := http.Server{
		Addr:    serverURL,
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		return err
	} else {
		log.Println("Server is stopped!")
	}

	return err
}

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}
	var h handlers.URLHandler

	h.Repo, err = repositories.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h.Auth = auth.NewAuthConfig(cfg.SecretKey)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Up!")

	router := router.NewRouter(&h)
	log.Println("Setup configuration!")

	err = StartServer(cfg.ServerURL, router)

	if err != nil {
		log.Fatal(err)
	}
}
