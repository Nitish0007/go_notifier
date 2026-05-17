package dto

import notifierif "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"

type SMTPConfiguration struct {
	Host     string `json:"host" validate:"required,min=1"`
	Port     int64  `json:"port" validate:"required,gt=0,lte=65535"`
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
	From     string `json:"from" validate:"required,email"`
}

func (c SMTPConfiguration) Channel() notifierif.Channel {
	return notifierif.ChannelEmail
}
