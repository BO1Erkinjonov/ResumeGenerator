package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type EmailData struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

func SendEmail(to []string, subject string, htmlpath string, code string) error {
	t, err := template.ParseFiles(htmlpath)
	if err != nil {
		log.Println(err)
		return err
	}

	var k bytes.Buffer
	err = t.Execute(&k, code)
	if err != nil {
		return err
	}
	if k.String() == "" {
		fmt.Println("Error buffer")
	}
	//mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//msg := []byte(fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, k.String()))
	// Authentication.
	auth := smtp.PlainAuth("", "boburerkinzonov@gmail.com", "llqmgbilccvhltfd", "smtp.gmail.com")

	// Sending email.
	err = smtp.SendMail("smtp.gmail.com:587", auth, "boburerkinzonov@gmail.com", to, []byte(code))
	return err
}
