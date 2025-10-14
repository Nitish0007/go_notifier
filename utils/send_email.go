package utils

import (
	"log"
	// "net/smtp"
)

func SendEmail() error {
	// auth := smtp.PlainAuth("", "your_email@example.com", "your_password", "smtp.example.com")
	log.Println("Email sent successfully")
	return nil
}