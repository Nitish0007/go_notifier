package utils

import (
	"log"
	"net/smtp"
	"fmt"
	"bytes"
	"crypto/tls"
	"github.com/google/uuid"

	"github.com/Nitish0007/go_notifier/internal/models"
)

func SendEmail(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody string, smtpConfig *models.SMTPConfiguration) error {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	// build email message
	message := buildEmailMessage(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody)

	// connect to smtp server
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("Error in connecting to SMTP server: %v", err)
		return err
	}
	defer client.Close()

	// start TLS
	tlsConfig := &tls.Config{
		ServerName: smtpConfig.Host,
		InsecureSkipVerify: true, // NOTE: can be false for development environment
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("Error in starting TLS: %v", err)
		return err
	}

	// set auth
  if err = client.Auth(auth); err != nil {
		log.Printf("Error in setting auth: %v", err)
		return err
	}

	// set sender
	if err = client.Mail(from); err != nil {
		log.Printf("Error in setting sender: %v", err)
		return err
	}

	// set recipient
	if err = client.Rcpt(to); err != nil {
		log.Printf("Error in setting recipient: %v", err)
		return err
	}

	// send email body
	writer, err := client.Data()
	if err != nil {
		log.Printf("Error in getting data writer: %v", err)
		return err
	}

	// write email body
	_, err = writer.Write(message)
	if err != nil {
		log.Printf("Error in writing email body: %v", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Printf("Error in closing data writer: %v", err)
		return err
	}

	// quit client
	err = client.Quit()
	if err != nil {
		log.Printf("Error in quitting client: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func buildEmailMessage(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody string) []byte {
	message := bytes.NewBuffer(nil)

	// set headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	message.WriteString(fmt.Sprintf("To: %s <%s>\r\n", toName, to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	
	message.WriteString(fmt.Sprintf("Content-Disposition: inline\r\n"))
	message.WriteString("\r\n")

	// add reply-to header
	if replyToEmail != "" && replyToName != "" {
		message.WriteString(fmt.Sprintf("Reply-To: %s <%s>\r\n", replyToName, replyToEmail))
	}
	message.WriteString("\r\n")

	// adding body to the message
	if htmlBody != "" {
		boundary := uuid.New().String()
		message.WriteString(fmt.Sprintf("Content-Type: multipart/alternative;\r\n boundary=\"%s\"\r\n", boundary))
		message.WriteString("\r\n")
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(body)
		message.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(htmlBody)
		message.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	}else{
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		message.WriteString("\r\n")
		message.WriteString(body)
	}

	return message.Bytes()
}