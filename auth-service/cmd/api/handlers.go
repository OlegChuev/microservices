package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := app.SuccessJson(
		fmt.Sprintf("Logged in user %s", user.Email),
		user,
	)

	app.WriteJson(w, http.StatusAccepted, payload)
}
