package content

import (
	"context"
	"errors"
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

func (r *ContentRepository) Index(ctx context.Context, accountID int64) ([]*Content, error) {
	var contents []*Content
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Find(&contents).Error
	if err != nil {
		return nil, errors.New("failed to get contents by account ID: " + err.Error())
	}

	return contents, nil
}

func (r *ContentRepository) GetById(ctx context.Context, accountID, id int64) (*Content, error) {
	var c Content
	err := r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContentRepository) Update(ctx context.Context, content *Content) error {
	return r.DB.WithContext(ctx).Save(content).Error
}

func (r *ContentRepository) Delete(ctx context.Context, accountID, id int64) error {
	return r.DB.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).Delete(&Content{}).Error
}