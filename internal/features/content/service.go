package content

import (
	"context"
)

type ContentService struct {
	contentRepository *ContentRepository
}

func NewContentService(contentRepository *ContentRepository) *ContentService {
	return &ContentService{
		contentRepository: contentRepository,
	}
}

func (s *ContentService) CreateContent(ctx context.Context, payload *CreateContentRequest) (*CreateContentResponse, error) {
	content := NewContent(payload.Content.AccountID, payload.Content.Body)
	err := s.contentRepository.Create(ctx, content)
	if err != nil {
		return nil, err
	}
	return &CreateContentResponse{
		ID: content.ID,
		AccountID: content.AccountID,
		Body: content.Body,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}, nil
}