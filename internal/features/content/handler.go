package content

import (
	"net/http"
	"strconv"

	"github.com/Nitish0007/go_notifier/internal/shared/api"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type ContentHandler struct {
	contentService *ContentService
}

func NewContentHandler(s *ContentService) *ContentHandler {
	return &ContentHandler{
		contentService: s,
	}
}

func (h *ContentHandler) GetContentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	contents, err := h.contentService.GetContents(ctx, accID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusOK, contents, "Contents fetched successfully")
}

func (h *ContentHandler) CreateContentHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[CreateContentRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	payload.Content.AccountID = accID

	content, err := h.contentService.CreateContent(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, content, "Content created successfully")
}

func (h *ContentHandler) GetContentByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	id, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid content ID")
		return
	}
	content, err := h.contentService.GetContentByID(ctx, accID, idInt64)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, content, "Content fetched successfully")
}

func (h *ContentHandler) UpdateContentHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[UpdateContentRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	id, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid content ID")
		return
	}
	payload.Content.ID = idInt64
	payload.Content.AccountID = accID
	content, err := h.contentService.UpdateContent(ctx, accID, idInt64, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, content, "Content updated successfully")
}

func (h *ContentHandler) DeleteContentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	id, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid content ID")
		return
	}
	err = h.contentService.DeleteContent(ctx, accID, idInt64)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, nil, "Content deleted successfully")
}