package list

import (
	"context"

	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/listsubscription"
	"gorm.io/gorm"
)

type ListRepository struct {
	db                   *gorm.DB
	listSubsRepo         *listsubscription.ListSubscriptionRepository
	contactRepo          *contact.ContactRepository
}

func NewListRepository(db *gorm.DB, lsr *listsubscription.ListSubscriptionRepository, cr *contact.ContactRepository	) *ListRepository {
	return &ListRepository{
		db: 					db,
		listSubsRepo: lsr,
		contactRepo:  cr,
	}
}

func (r *ListRepository) Create(ctx context.Context, list *List) error {
	return r.db.WithContext(ctx).Create(list).Error
}

func (r *ListRepository) Index(ctx context.Context, accID int64) ([]*List, error) {
	var lists []*List
	err := r.db.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&lists).Error
	if err != nil {
		return nil, err
	}
	
	return lists, nil
}

func (r *ListRepository) FindContactByUUID(ctx context.Context,accId int64, uuid string) (*contact.Contact, error) {
	contact, err := r.contactRepo.FindByUUID(ctx, accId, uuid)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ListRepository) FindById(ctx context.Context, accId int64, id int64) (*contact.Contact, error) {
	contact, err := r.contactRepo.FindById(ctx, accId, id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ListRepository) FindOrCreateEmailContact(ctx context.Context, accId int64, contactPayload *contact.ContactPayload) (*contact.Contact, error) {
	contact, err := r.contactRepo.FindOrCreateByEmail(ctx, accId, contactPayload)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ListRepository) CreateListSubscription(ctx context.Context, reqData *SubscribeToListRequest) error {
	listSubscription := listsubscription.NewListSubscription(reqData.AccountID, reqData.ListID, reqData.ContactID, reqData.Active)
	return r.listSubsRepo.Create(ctx, listSubscription)
}