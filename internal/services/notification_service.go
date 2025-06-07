package services

import (
	"log"

	"github.com/Nitish0007/go_notifier/internal/notifiers"
)

type NotificationService struct {
	notifiers map[string]notifiers.Notifier
}

func NewNotificationService(list []notifiers.Notifier) *NotificationService {
	nList := make(map[string]notifiers.Notifier)
	for _, val := range list {
		if val.ChannelType() == "email" || val.ChannelType() == "sms" || val.ChannelType() == "in_app" {
			nList[val.ChannelType()] = val
		}else{
			log.Printf(">>>>>>>>>>>>> Unknown Notifier: %v", val.ChannelType())
		}
	}
	return &NotificationService{
		notifiers: nList,
	}
}

