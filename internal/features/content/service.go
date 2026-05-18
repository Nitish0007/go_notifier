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
	content := NewContent(payload.Content.AccountID, payload.Content.Name, payload.Content.Body)
	err := s.contentRepository.Create(ctx, content)
	if err != nil {
		return nil, err
	}
	return &CreateContentResponse{
		ID: content.ID,
		Name: content.Name,
		AccountID: content.AccountID,
		Body: content.Body,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}, nil
}

func (s *ContentService) GetContentByID(ctx context.Context, accID int64, id int64) (*ContentResponse, error) {
	content, err := s.contentRepository.GetByID(ctx, accID, id)
	if err != nil {
		return nil, err
	}
	
	return &ContentResponse{
		ID: content.ID,
		Name: content.Name,
		AccountID: content.AccountID,
		Body: content.Body,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}, nil
}

func (s *ContentService) GetContents(ctx context.Context, accID int64) ([]*ContentResponse, error) {
	contents, err := s.contentRepository.Index(ctx, accID)
	if err != nil {
		return nil, err
	}
	contentResponses := make([]*ContentResponse, len(contents))
	for i, content := range contents {
		contentResponses[i] = &ContentResponse{
			ID: content.ID,
			Name: content.Name,
			AccountID: content.AccountID,
			Body: content.Body,
			CreatedAt: content.CreatedAt,
			UpdatedAt: content.UpdatedAt,
		}
	}
	return contentResponses, nil
}

func (s *ContentService) UpdateContent(ctx context.Context, accID int64, id int64, payload *UpdateContentRequest) (*ContentResponse, error) {
	content, err := s.contentRepository.GetById(ctx, accID, id)
	if err != nil {
		return nil, err
	}
	content.ID = payload.Content.ID
	content.AccountID = payload.Content.AccountID
	content.Name = payload.Content.Name
	content.Body = payload.Content.Body
	err = s.contentRepository.Update(ctx, content)
	if err != nil {
		return nil, err
	}
	return &ContentResponse{
		ID: content.ID,
		Name: content.Name,
		AccountID: content.AccountID,
		Body: content.Body,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}, nil
}

func (s *ContentService) DeleteContent(ctx context.Context, accID int64, id int64) error {
	err := s.contentRepository.Delete(ctx, accID, id)
	if err != nil {
		return err
	}
	return nil
}