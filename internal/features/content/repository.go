package content

import (
	"context"
	"gorm.io/gorm"
)

type ContentRepository struct {
	DB *gorm.DB
}

func NewContentRepository(conn *gorm.DB) *ContentRepository {
	return &ContentRepository{
		DB: conn,
	}
}

func (r *ContentRepository) Create(ctx context.Context, content *Content) error {
	return r.DB.WithContext(ctx).Create(content).Error
}