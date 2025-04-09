package main

import (
	"fmt"
	"go-svc-gophermart/internal/config"
	"go-svc-gophermart/internal/handlers"
	"go-svc-gophermart/internal/routerAPI"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	//var t repositories.RepoDBPostgres

	//repositories.Repository.NewRepository(cfg.DSN)// .CreateStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	/*
		auth := auth.NewAuthConfig()
	*/
	log.Println("Up!")
	var h handlers.URLHandler

	router := routerAPI.NewRouter(&h)
	log.Println("Setup configuration!")

	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
