package handlers

import (
	"net/http"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/services"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(s *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: s,
	}
}

func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ParseJSONBody(r)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	accountData, exists := payload["account"].(map[string]any)
	if !exists || len(accountData) == 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Account data is required")
		return
	}

	// Initialize account with the provided data
	ctx := r.Context()
	account, err := h.accountService.InitializeAccount(ctx, accountData)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// account creation and related operations handled here
	account, err = h.accountService.CreateAccount(ctx, account)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, account, "Account created successfully")
}