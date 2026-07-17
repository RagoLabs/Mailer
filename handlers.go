package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

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

