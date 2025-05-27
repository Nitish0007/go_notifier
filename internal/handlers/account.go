package handlers

import (
	"net/http"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/services"
)

var accountService = services.NewAccountService()

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ParseJSONBody(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountData, exists := payload["account"].(map[string]any)
	if !exists || len(accountData) == 0 {
		http.Error(w, "Invalid account data", http.StatusBadRequest)
		return
	}

	// Initialize account with the provided data
	ctx := r.Context()
	account, err := accountService.InitializeAccount(ctx, accountData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// account creation and related operations handled here
	account, err = accountService.CreateAccount(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, account)
}