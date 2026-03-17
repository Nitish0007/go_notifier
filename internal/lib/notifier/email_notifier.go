package notifier

import (
	"fmt"
	"log"
	"context"
	"errors"
	"net/smtp"
	"crypto/tls"
	"bytes"
	"github.com/google/uuid"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/shared/dto"
	"github.com/Nitish0007/go_notifier/internal/features/notification"
	notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

type EmailNotifier struct {
	notificationRepository *notification.NotificationRepository
};

func NewEmailNotifier(r *notification.NotificationRepository) *EmailNotifier {
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
	notif := &notification.Notification{}
	aid := utils.GetCurrentAccountID(ctx)
	if aid == -1 {
		return nil, errors.New("account ID is required to create a notification")
	}
	notif.AccountID = aid
	emailChannel, exists := payload["channel"].(string)
	if !exists || emailChannel != "email" {
		return nil, errors.New("channel is required and must be 'email'")
	}

	subject, exists := payload["subject"].(string)
	if !exists || subject == "" {
		return nil, errors.New("subject is required to create notification")
	}
	notif.Subject = subject

	nChannel, err := notification.StringToNotificationChannel(emailChannel)
	if err != nil {
		return nil, err
	}

	// NOTE: THE PAYLOAD STRUCTURE MUST BE REFINED LATER, THIS IS JUST A PROTOTYPE
	notif.Channel = nChannel
	nStatus, err := notification.StringToNotificationStatus("pending")
	if err != nil {
		return nil, err
	}
	notif.Status = nStatus
	notif.Recipient = payload["recipient"].(string)
	notif.Body = payload["body"].(string)
	notif.HtmlBody = payload["html_body"].(string)
	sendTime := payload["send_at"]
	if sendTime != nil {
		notif.SendAt = utils.ParseTime(sendTime.(string))
	} else {
		notif.SendAt = utils.GetCurrentTime()
	}
	
	sanitizedMetadata := make(map[string]any)	
	if metadata, exists := payload["metadata"].(map[string]any); exists {
		if len(metadata) == 0 {
			return nil, errors.New("metadata can't be empty")
		}
		if _, ok := metadata["from_email"]; !ok || metadata["from_email"] == "" {
			return nil, errors.New("from_email is required in metadata")
		}
		if _, ok := metadata["from_name"]; !ok || metadata["from_name"] == "" {
			return nil, errors.New("from_name is required in metadata")
		}
		if _, ok := metadata["to_name"]; !ok || metadata["to_name"] == "" {
			return nil, errors.New("to_name is required in metadata")
		}
		if _, ok := metadata["reply_to_email"]; !ok || metadata["reply_to_email"] == "" {
			return nil, errors.New("reply_to_email is required in metadata")
		}
		if _, ok := metadata["reply_to_name"]; !ok || metadata["reply_to_name"] == "" {
			return nil, errors.New("reply_to_name is required in metadata")
		}
		

		// keeping only required fields in metadata
		sanitizedMetadata["from_email"] = metadata["from_email"].(string)
		sanitizedMetadata["from_name"] = metadata["from_name"].(string)
		sanitizedMetadata["to_name"] = metadata["to_name"].(string)
		sanitizedMetadata["reply_to_email"] = metadata["reply_to_email"].(string)
		sanitizedMetadata["reply_to_name"] = metadata["reply_to_name"].(string)

		notif.Metadata = sanitizedMetadata
		// notification.ErrorMessage = nil
	} else {
		return nil, errors.New("metadata is required in notification data")
	}

	err = n.notificationRepository.Create(ctx, notif)
	return notif, err
}

func (n *EmailNotifier) CreateBulkNotifications(ctx context.Context, payload []map[string]any) ([]*notification.Notification, error) {
	// Initialize the notification model instance for given payload
	// notifications := []*models.Notification{}
	// for _, p := range payload {
	// 	notification := &models.Notification{}
	// 	notifications = append(notifications, notification)
	// }

	// err := n.notificationRepository.Create(ctx, notifications)
	// return notifications, err
	return nil, nil
}

func (n *EmailNotifier) ChannelType() string {
	return "email"
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