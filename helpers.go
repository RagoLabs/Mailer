package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

type ContactForm struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	Service   string `json:"service" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

type SignupData struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type ActivateOrResetData struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type ResetCompleteData struct {
	Email string `json:"email" validate:"required,email"`
}


func (app *application) sendContactUsEmail(form ContactForm, recipients []string) error {

	type templateData struct {
		ContactForm
		FormattedDate string
	}

	data := templateData{
		ContactForm:   form,
		FormattedDate: time.Now().Format("January 2, 2006 at 3:04 PM"),
	}

	tmpl, err := template.New("email").Parse(contactusTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	toHeader := ""
	for i, recipient := range recipients {
		if i > 0 {
			toHeader += ", "
		}
		toHeader += recipient
	}

	subject := "Contact Form Submission"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, toHeader, emailBody.String())

	auth := smtp.PlainAuth("", app.config.mail.user, app.config.mail.pwd, app.config.mail.host)
	err = smtp.SendMail(app.config.mail.host+":"+app.config.mail.port, auth, app.config.mail.user, recipients, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}



