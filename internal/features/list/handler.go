package list

import (
	"net/http"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type ListHandler struct {
	listService *ListService
}

func NewListHandler(s *ListService) *ListHandler {
	return &ListHandler{
		listService: s,
	}
}

func (h *ListHandler) CreateListHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[CreateListRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	list, err := h.listService.CreateList(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, list, "List created successfully")
}

func (h *ListHandler) GetListsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	lists, err := h.listService.GetLists(ctx, accID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, lists, "Lists fetched successfully")
}

func (h *ListHandler) SubscribeToListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	listID, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx = sharedhelper.SetValueToContext(ctx, "listID", listID)

	payload, err := api.ParseAndValidateRequest[SubscribeToListRawPayload](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.listService.SubscribeToList(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, response, "Contacts subscribed to list successfully")

}

func (h *ListHandler) ManageListSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// accID := sharedhelper.GetCurrentAccountID(ctx)
	// listID := chi.URLParam(r, "id")
	// payload, err := api.ParseAndValidateRequest[ManageListSubscriptionRequest](r)
	// if err != nil {
	// 	api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
}