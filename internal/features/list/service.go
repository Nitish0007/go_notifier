package list

import (
	"fmt"
	"time"
	"context"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
)

type ListService struct {
	listRepository *ListRepository
}

func NewListService(listRepository *ListRepository) *ListService {
	return &ListService {
		listRepository: listRepository,
	}
}

func (s *ListService) GetLists(ctx context.Context, accID int64) ([]*CreateListResponse, error) {
	lists, err := s.listRepository.Index(ctx, accID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}

	listResponses := make([]*CreateListResponse, len(lists))
	for i, list := range lists {
		listResponses[i] = &CreateListResponse{
			ID: list.ID,
			AccountID: list.AccountID,
			Name: list.Name,
			ContactsCount: list.ContactsCount,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		}
	}
	return listResponses, nil
}

func (s *ListService) CreateList(ctx context.Context, payload *CreateListRequest) (*CreateListResponse, error) {
	list := NewList(payload.List.AccountID, payload.List.Name)
	err := s.listRepository.Create(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}
	return &CreateListResponse{
		ID: list.ID,
		AccountID: list.AccountID,
		Name: list.Name,
		ContactsCount: list.ContactsCount,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}, nil
}

func (s *ListService) SubscribeToList(ctx context.Context, payload *SubscribeToListRawPayload) (*SubscribeToListResponse, error) {
	listID := sharedhelper.GetValueFromContext(ctx, "listID")
	if listID == nil {
		return nil, fmt.Errorf("list id not found in context")
	}
	listIDInt, ok := listID.(int64)
	if !ok {
		return nil, fmt.Errorf("list id is not an integer")
	}

	accId := int64(sharedhelper.GetCurrentAccountID(ctx))

	var reqData *SubscribeToListRequest
	reqData = &SubscribeToListRequest{
		AccountID: accId,
		ListID: listIDInt,
		Active: payload.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if payload.ContactID > 0 {
		contact, err := s.listRepository.FindById(ctx, accId, payload.ContactID)
		if err != nil {
			return nil, fmt.Errorf("failed to find contact by id: %w", err)
		}
		reqData.ContactID = contact.ID
	} else if payload.UUID != "" {
		contact, err := s.listRepository.FindContactByUUID(ctx, accId, payload.UUID)
		if err != nil {
			return nil, fmt.Errorf("failed to find contact by uuid: %w", err)
		}
		reqData.ContactID = contact.ID
	} else if payload.EmailContact.Email != "" {
		contact, err := s.listRepository.FindOrCreateEmailContact(ctx, accId, &contact.ContactPayload{
			Email: payload.EmailContact.Email,
			FirstName: payload.EmailContact.FirstName,
			LastName: payload.EmailContact.LastName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find or create contact by email: %w", err)
		}
		reqData.ContactID = contact.ID
	}

	err := s.listRepository.CreateListSubscription(ctx, reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to create list subscription: %w", err)
	}
	return &SubscribeToListResponse{
		ListID: reqData.ListID,
		ContactID: reqData.ContactID,
		Active: reqData.Active,
		CreatedAt: reqData.CreatedAt,
		UpdatedAt: reqData.UpdatedAt,
	}, nil
}