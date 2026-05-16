package list

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/listsubscription"
	"gorm.io/gorm"
)

type ListRepository struct {
	DB           *gorm.DB
	listSubsRepo *listsubscription.ListSubscriptionRepository
	contactRepo  *contact.ContactRepository
}

func NewListRepository(db *gorm.DB, lsr *listsubscription.ListSubscriptionRepository, cr *contact.ContactRepository) *ListRepository {
	return &ListRepository{
		DB:           db,
		listSubsRepo: lsr,
		contactRepo:  cr,
	}
}

func (r *ListRepository) Create(ctx context.Context, list *List) error {
	return r.DB.WithContext(ctx).Create(list).Error
}

func (r *ListRepository) Index(ctx context.Context, accID int64) ([]*List, error) {
	var lists []*List
	err := r.DB.WithContext(ctx).Where("account_id = ?", accID).Order("created_at DESC").Find(&lists).Error
	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *ListRepository) FindContactByUUID(ctx context.Context, accId int64, uuid string) (*contact.Contact, error) {
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

// SubscribeToList resolves the contact (find or create by email), verifies the list belongs to the account,
// and inserts the list subscription in a single transaction.
func (r *ListRepository) SubscribeToList(ctx context.Context, accountID, listID int64, payload *SubscribeToListRawPayload) (*SubscribeToListResponse, error) {
	var out SubscribeToListResponse
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var listRow List
		if err := tx.Where("id = ? AND account_id = ?", listID, accountID).First(&listRow).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("list not found")
			}
			return err
		}

		var cont *contact.Contact
		var err error

		switch {
		case payload.ContactID > 0:
			cont, err = r.contactRepo.FindByIDWithTx(ctx, tx, accountID, payload.ContactID)
			if err != nil {
				return fmt.Errorf("failed to find contact by id: %w", err)
			}
		case payload.UUID != "":
			cont, err = r.contactRepo.FindByUUIDWithTx(ctx, tx, accountID, payload.UUID)
			if err != nil {
				return fmt.Errorf("failed to find contact by uuid: %w", err)
			}
		case payload.EmailContact.Email != "":
			cont, err = r.contactRepo.FindOrCreateByEmailWithTx(ctx, tx, accountID, &contact.ContactPayload{
				Email:     payload.EmailContact.Email,
				FirstName: payload.EmailContact.FirstName,
				LastName:  payload.EmailContact.LastName,
			})
			if err != nil {
				return fmt.Errorf("failed to find or create contact by email: %w", err)
			}
		default:
			return errors.New("contact_id, uuid, or email_contact.email is required")
		}

		sub := listsubscription.NewListSubscription(accountID, listID, cont.ID, payload.Active)
		if err := r.listSubsRepo.CreateWithTx(ctx, tx, sub); err != nil {
			return fmt.Errorf("failed to create list subscription: %w", err)
		}

		out = SubscribeToListResponse{
			ListID:      listID,
			ContactID:   cont.ID,
			ContactUUID: cont.UUID,
			Active:      sub.Active,
			CreatedAt:   sub.CreatedAt,
			UpdatedAt:   sub.UpdatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ListRepository) GetSubscribers(ctx context.Context, accountId int64, listID int64) ([]*SubscriberResponse, error) {
	var subscribers = make([]*SubscriberResponse, 0)
	err := r.DB.WithContext(ctx).Raw(`
		SELECT c.id, c.first_name, c.last_name, c.uuid, c.created_at as contact_created_at, c.updated_at as contact_updated_at, ls.active, ls.created_at as subscription_created_at, ls.updated_at as subscription_updated_at
		FROM list_subscriptions as ls
		INNER JOIN contacts as c
		ON ls.contact_id = c.id
		WHERE ls.list_id = ? AND ls.account_id = ? AND c.account_id = ?
	`, listID, accountId, accountId).Scan(&subscribers).Error

	if err != nil {
		return nil, err
	}

	return subscribers, nil
}
