package content

import (
	"time"
	"errors"
	"strings"

	"gorm.io/gorm"
	"github.com/flosch/pongo2/v7"
)

type Content struct {
	ID        int64      `json:"id" gorm:"primaryKey" validate:"omitempty,gt=0"`
	Name      string     `json:"name" gorm:"not null" validate:"required"`
	AccountID int64      `json:"account_id" gorm:"not null;index" validate:"required,gt=0"`
	Body      string     `json:"body" gorm:"not null" validate:"required"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`
}

func NewContent(accountId int64, name, body string) *Content {
	return &Content{
		AccountID: accountId,
		Name: name,
		Body: body,
	}
}

func (c *Content) BeforeSave(tx *gorm.DB) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if c.Body == "" {
		return errors.New("body is required")
	}
	_, err := pongo2.FromString(c.Body)
	return err
}