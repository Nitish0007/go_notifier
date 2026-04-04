package notifier

import (
	"fmt"
	"log"
	"context"
	// "errors"
	"net/smtp"
	"crypto/tls"
	"bytes"
	"github.com/google/uuid"


	"github.com/Nitish0007/go_notifier/internal/shared/dto"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	// "github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

type EmailNotifier struct {
	notificationRepository *emailnotification.EmailNotificationRepository
};

func NewEmailNotifier(r *emailnotification.EmailNotificationRepository) *EmailNotifier {
	return &EmailNotifier{
		notificationRepository: r,
	}
}

func (n *EmailNotifier) Notify(notificationView notifierInterface.NotificationView, smtpConfig *dto.SMTPConfiguration) error {
	from := notificationView.GetMetadata()["from_email"].(string)
	fromName := notificationView.GetMetadata()["from_name"].(string)
	to := notificationView.GetRecipient()
	toName := notificationView.GetMetadata()["to_name"].(string)
	replyToEmail := notificationView.GetMetadata()["reply_to_email"].(string)
	replyToName := notificationView.GetMetadata()["reply_to_name"].(string)
	subject := notificationView.GetSubject()
	body := notificationView.GetBody()
	htmlBody := notificationView.GetHtmlBody()
	err := Send(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody, smtpConfig)
	if err != nil {
		return err
	}
	return nil
}

func (n *EmailNotifier) CreateNotification(ctx context.Context, payload map[string]any) (notifierInterface.NotificationView, error) {
	// Initialize the notification model instance for given payload
	// notif := &emailnotification.EmailNotification{}
	// aid := sharedhelper.GetCurrentAccountID(ctx)
	// if aid == -1 {
	// 	return nil, errors.New("account ID is required to create a notification")
	// }
	// notif.AccountID = aid
	// emailChannel, exists := payload["channel"].(string)
	// if !exists || emailChannel != "email" {
	// 	return nil, errors.New("channel is required and must be 'email'")
	// }

	// subject, exists := payload["subject"].(string)
	// if !exists || subject == "" {
	// 	return nil, errors.New("subject is required to create notification")
	// }
	// notif.Subject = subject

	

	// // NOTE: THE PAYLOAD STRUCTURE MUST BE REFINED LATER, THIS IS JUST A PROTOTYPE
	// notif.Status, err := emailnotification.StringToEmailNotificationStatus(payload["status"].(string))
	// if err != nil {
	// 	return nil, err
	// }	


	// if err := n.notificationRepository.Create(ctx, notif); err != nil {
	// 	return nil, err
	// }
	// return notif, nil
	return nil, nil
}


func Send(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody string, smtpConfig *dto.SMTPConfiguration) error {
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
	
	// to check if the client is closed before sending the email
	defer func() {
		log.Printf("Closing SMTP client")
		if err := client.Close(); err != nil {
			log.Printf("Warning: Error in closing SMTP client: %v", err)
		}
	}()

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
		writer.Close()
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
		// log.Printf("Error in quitting client: %v", err)
		log.Printf("Warning: Error in quitting client: %v", err)
		// return nil
	}

	log.Println("Email sent successfully")
	return nil
}

func buildEmailMessage(to, from, fromName, toName, replyToEmail, replyToName, subject, body, htmlBody string) []byte {
	message := bytes.NewBuffer(nil)

	// set headers
	message.WriteString(fmt	.Sprintf("From: %s <%s>\r\n", fromName, from)) // NOTE: %s is used to format the string and \r\n is used to add a new line
	message.WriteString(fmt.Sprintf("To: %s <%s>\r\n", toName, to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	
	message.WriteString("Content-Disposition: inline\r\n")
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