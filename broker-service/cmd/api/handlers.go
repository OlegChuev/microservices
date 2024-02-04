package main

import (
	"net/http"

	"github.com/OlegChuev/microservices/utils"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := utils.JsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	utils.WriteJson(w, http.StatusOK, payload)
}
