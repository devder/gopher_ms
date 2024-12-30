package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Reuse the HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// send json to microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	req, err := http.NewRequest(http.MethodPost, "http://auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer resp.Body.Close()

	// ensure we get back the correct status
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode != http.StatusAccepted {
		errMsg := "failed to call auth service"
		if resp.StatusCode == http.StatusUnauthorized {
			errMsg = "invalid credentials"
		}
		app.errorJSON(w, errors.New(errMsg))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	// send json to microservice
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	// call the service
	req, err := http.NewRequest(http.MethodPost, "http://logger/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer resp.Body.Close()

	// ensure we get back the correct status
	if resp.StatusCode != http.StatusAccepted {
		errMsg := "failed to call logger service"
		app.errorJSON(w, errors.New(errMsg))
		return
	}

	// create a variable we'll read response.Body into
	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// call the mail service
	req, err := http.NewRequest(http.MethodPost, "http://mailer/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer resp.Body.Close()

	// ensure we get back the correct status
	if resp.StatusCode != http.StatusAccepted {
		errMsg := "failed to call mailer service"
		app.errorJSON(w, errors.New(errMsg))
		return
	}

	// create a variable we'll read response.Body into
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("mail sent to %s", m.To),
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
