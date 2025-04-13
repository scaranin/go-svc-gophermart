package main

import (
	"go-svc-gophermart/internal/auth"
	"go-svc-gophermart/internal/config"
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/repositories"
	"go-svc-gophermart/internal/routerapi"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	var h handlers.URLHandler

	h.Repo, err = repositories.NewRepository(cfg.DSN)

	h.Auth = auth.NewAuthConfig()

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Up!")

	router := routerapi.NewRouter(&h)
	log.Println("Setup configuration!")

	log.Println("Start server on ", cfg.ServerURL)
	err = http.ListenAndServe(cfg.ServerURL, router)

	if err != nil {
		log.Fatal(err)
	}
}
