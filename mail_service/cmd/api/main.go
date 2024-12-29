package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "80"

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return Mail{
		domain:      os.Getenv("SMTP_DOMAIN"),
		host:        os.Getenv("SMTP_HOST"),
		port:        port,
		username:    os.Getenv("SMTP_USER"),
		password:    os.Getenv("SMTP_PASSWORD"),
		encryption:  os.Getenv("SMTP_ENCRYPTION"),
		fromAddress: os.Getenv("EMAIL_NOREPLY_ADD"),
		fromName:    os.Getenv("EMAIL_NOREPLY_NAME"),
	}

}
