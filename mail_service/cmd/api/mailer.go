package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	domain      string
	host        string
	port        int
	username    string
	password    string
	encryption  string
	fromAddress string
	fromName    string
}

type Message struct {
	from        string
	fromName    string
	to          string
	subject     string
	attachments []string
	data        any
	dataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.from == "" {
		msg.from = m.fromAddress
	}

	if msg.fromName == "" {
		msg.fromName = m.fromName
	}

	data := map[string]any{
		"message": msg.data,
	}

	msg.dataMap = data
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.host
	server.Port = m.port
	server.Username = m.username
	server.Password = m.password
	server.Encryption = m.getEncryption(m.encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.from).
		AddTo(msg.to).
		SetSubject(msg.subject).
		SetBody(mail.TextPlain, plainMessage).
		AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.attachments) > 0 {
		for _, x := range msg.attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.dataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.dataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	return formattedMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &opts)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
