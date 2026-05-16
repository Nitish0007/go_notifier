package listsubscription

import (
	"context"
	"gorm.io/gorm"
)

type ListSubscriptionRepository struct {
	db *gorm.DB
}

func NewListSubscriptionRepository(db *gorm.DB) *ListSubscriptionRepository {
	return &ListSubscriptionRepository{db: db}
}

func (r *ListSubscriptionRepository) Create(ctx context.Context, listSubscription *ListSubscription) error {
	return r.db.WithContext(ctx).Create(listSubscription).Error
}

func (r *ListSubscriptionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, listSubscription *ListSubscription) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(listSubscription).Error
}
