package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"github.com/devder/gopher_ms/broker/event"
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

type RPCPayload struct {
	Name string
	Data string
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
		// app.logItem(w, requestPayload.Log)
		// app.logEventViaRabbit(w, requestPayload.Log)
		app.logEventViaRPC(w, requestPayload.Log)
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

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged via RabbitMQ",
		Data:    l.Name,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload(l)

	var result string
	// call the service with the exact method name (LogInfo) as defined in the RPCServer struct on the logger service
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter := event.NewEventEmitter(app.rabbitConn)

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	p, _ := json.Marshal(&payload)
	return emitter.Push(string(p), "log.INFO")
}
