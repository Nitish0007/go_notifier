package notifier

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	// "github.com/google/uuid"

	"github.com/Nitish0007/go_notifier/internal/shared/dto"
	notifierif "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

// EmailNotifier delivers messages over SMTP.
type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (n *EmailNotifier) Channel() notifierif.Channel {
	return notifierif.ChannelEmail
}

func (n *EmailNotifier) Notify(ctx context.Context, req notifierif.DeliveryRequest, cfg notifierif.ProviderConfig) error {
	_ = ctx
	smtpCfg, ok := cfg.(*dto.SMTPConfiguration)
	if !ok || smtpCfg == nil {
		return fmt.Errorf("email notifier: expected *dto.SMTPConfiguration, got %T", cfg)
	}

	from := req.Metadata["from_email"]
	if from == "" {
		from = smtpCfg.From
	}
	fromName := req.Metadata["from_name"]
	toName := req.Metadata["to_name"]
	replyToEmail := req.Metadata["reply_to_email"]
	replyToName := req.Metadata["reply_to_name"]

	body := req.Body

	return sendSMTP(req.Recipient, from, fromName, toName, replyToEmail, replyToName, req.Subject, body, smtpCfg)
}

func sendSMTP(to, from, fromName, toName, replyToEmail, replyToName, subject, body string, smtpConfig *dto.SMTPConfiguration) error {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	message := buildEmailMessage(to, from, fromName, toName, replyToEmail, replyToName, subject, body)

	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("smtp dial: %v", err)
		return err
	}
	quitChan := false
	defer func() {
		if quitChan {
			return
		}

		if err := client.Close(); err != nil {
			log.Printf("smtp close: %v", err)
		}
		quitChan = true
	}()

	tlsConfig := &tls.Config{
		ServerName:         smtpConfig.Host,
		InsecureSkipVerify: true,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("smtp starttls: %v", err)
		return err
	}
	if err = client.Auth(auth); err != nil {
		log.Printf("smtp auth: %v", err)
		return err
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}
	if err = client.Quit(); err != nil {
		log.Printf("smtp quit: %v", err)
	}
	quitChan = true
	return nil
}

func buildEmailMessage(to, from, fromName, toName, replyToEmail, replyToName, subject, body string) []byte {
	message := bytes.NewBuffer(nil)
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	message.WriteString(fmt.Sprintf("To: %s <%s>\r\n", toName, to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Disposition: inline\r\n\r\n")

	if replyToEmail != "" && replyToName != "" {
		message.WriteString(fmt.Sprintf("Reply-To: %s <%s>\r\n", replyToName, replyToEmail))
	}
	message.WriteString("\r\n")

	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	message.WriteString(body)
	return message.Bytes()
}
