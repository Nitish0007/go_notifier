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

func (r *ContentRepository) GetByID(ctx context.Context, accountID, id int64) (*Content, error) {
	var c Content
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}