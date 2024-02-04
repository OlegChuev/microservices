package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OlegChuev/microservices/utils"
)

const WEB_PORT = "80"

type Config struct {
	*utils.Config
}

func main() {
	app := Config{}

	log.Printf("Starting Broker service on the %s port...\n", WEB_PORT)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
