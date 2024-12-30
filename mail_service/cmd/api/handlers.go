package main

import (
	"fmt"
	"net/http"
	"os"
)

type mailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		from:    os.Getenv("EMAIL_NOREPLY_ADD"),
		to:      requestPayload.To,
		subject: requestPayload.Subject,
		data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("sent to %s", requestPayload.To),
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
