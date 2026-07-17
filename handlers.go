package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ContactForm struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	Service   string `json:"service" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

func (app *application) sendContactEmailHandler(c echo.Context) error {

	var input ContactForm

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	recipients := strings.Split(app.config.mail.recipients, ",")

	if err := app.sendContactUsEmail(input, recipients); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send emails"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Emails sent successfully!"})
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
