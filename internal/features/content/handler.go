package content

import (
	"net/http"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
)

type ContentHandler struct {
	contentService *ContentService
}

func NewContentHandler(s *ContentService) *ContentHandler {
	return &ContentHandler{
		contentService: s,
	}
}

func (h *ContentHandler) CreateContentHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[CreateContentRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	content, err := h.contentService.CreateContent(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, content, "Content created successfully")
}