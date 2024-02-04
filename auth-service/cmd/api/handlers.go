package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OlegChuev/microservices/utils"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := utils.ReadJson(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		utils.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		utils.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	payload := utils.SuccessJson(fmt.Sprintf("Logged in user %s", user.Email))

	utils.WriteJson(w, http.StatusOK, payload)
}
