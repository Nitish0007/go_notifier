package list

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type ListService struct {
	listRepository *ListRepository
}

func NewListService(listRepository *ListRepository) *ListService {
	return &ListService{
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
			ID:            list.ID,
			AccountID:     list.AccountID,
			Name:          list.Name,
			ContactsCount: list.ContactsCount,
			CreatedAt:     list.CreatedAt,
			UpdatedAt:     list.UpdatedAt,
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
		ID:            list.ID,
		AccountID:     list.AccountID,
		Name:          list.Name,
		ContactsCount: list.ContactsCount,
		CreatedAt:     list.CreatedAt,
		UpdatedAt:     list.UpdatedAt,
	}, nil
}

func (s *ListService) SubscribeToList(ctx context.Context, payload *SubscribeToListRawPayload) (*SubscribeToListResponse, error) {
	listIDVal := sharedhelper.GetValueFromContext(ctx, "listID")
	if listIDVal == nil {
		return nil, fmt.Errorf("list id not found in context")
	}
	listIDStr, ok := listIDVal.(string)
	if !ok || listIDStr == "" {
		return nil, fmt.Errorf("invalid list id in context")
	}
	listIDInt, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list id: %w", err)
	}

	accID := sharedhelper.GetCurrentAccountID(ctx)
	return s.listRepository.SubscribeToList(ctx, accID, listIDInt, payload)
}

func (s *ListService) GetSubscribers(ctx context.Context) ([]*SubscriberResponse, error) {
	accId := sharedhelper.GetCurrentAccountID(ctx)
	listId := sharedhelper.GetValueFromContext(ctx, "listID")
	listIdInt64, err := strconv.ParseInt(listId.(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list id: %w", err)
	}

	subscribers, err := s.listRepository.GetSubscribers(ctx, accId, listIdInt64)
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}