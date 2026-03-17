package account

import (
	"log"
	"net/http"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
)

type AccountHandler struct {
	accountService *AccountService
}

func NewAccountHandler(s *AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: s,
	}
}

func (h *AccountHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := api.ParseAndValidateRequest[SignupRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	account, err := h.accountService.CreateAccount(ctx, reqData)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, account, "Account created successfully")
}

func (h *AccountHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := api.ParseAndValidateRequest[LoginRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	response, err := h.accountService.Login(ctx, reqData)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>> response", response)

	api.WriteResponse(w, http.StatusOK, response, "Login successful")
}