package contact

import (
	"net/http"
	"strconv"
	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
)

type ContactHandler struct {
	contactService    *ContactService
}

func NewContactHandler(s *ContactService) *ContactHandler {
	return &ContactHandler {
		contactService: s,
	}
}

func (h *ContactHandler) GetContactsHandler(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	accID := utils.GetCurrentAccountID(ctx)
	contacts, err := h.contactService.GetContacts(ctx, accID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, contacts, "Contacts fetched successfully")
}

func (h *ContactHandler) CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqData, err := api.ParseAndValidateRequest[CreateContactRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	contact, err := h.contactService.CreateContact(ctx, reqData)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusCreated, contact, "Contact created successfully")
}

func (h *ContactHandler) GetContactByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	ctx := r.Context()
	contact, err := h.contactService.GetContactByKey(ctx, "id", idInt)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, contact, "Contact fetched successfully")
}

func (h *ContactHandler) GetContactByUUIDHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := utils.GetPathParam(r, "uuid")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	contact, err := h.contactService.GetContactByKey(ctx, "uuid", uuid)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, contact, "Contact fetched successfully")
}