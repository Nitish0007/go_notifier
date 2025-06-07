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

// for signup
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
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, account, "Account created successfully")
}

// for login
func (h *AccountHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ParseJSONBody(r)

	if err != nil{
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	loginData, exists := payload["login"].(map[string]any)
	if(!exists || len(loginData) == 0){
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Credentials not provided")
		return
	}
	
	ctx := r.Context()
	apiKey, err := h.accountService.Login(ctx, loginData)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	data := make(map[string]string)
	data["auth_token"] = apiKey

	utils.WriteJSONResponse(w, http.StatusOK, data, "Logged in successfully")

}