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