package handlers

import (
	"net/http"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/services"
	"github.com/Nitish0007/go_notifier/internal/validators"
	"github.com/Nitish0007/go_notifier/utils"
)

type ConfigurationHandler struct {
	configurationService *services.ConfigurationService
}

func NewConfigurationHandler(s *services.ConfigurationService) *ConfigurationHandler {
	return &ConfigurationHandler{
		configurationService: s,
	}
}

// for getting configurations
func (h *ConfigurationHandler) GetConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := utils.GetCurrentAccountID(ctx)
	configs, err := h.configurationService.GetConfigurations(ctx, accID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, configs, "Configurations fetched successfully")
}

// for creating configuration
func (h *ConfigurationHandler) CreateConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ParseJSONBody(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	configData, exists := payload["configuration"].(map[string]any)
	if !exists || len(configData) == 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Configuration payload is required")
		return
	}

	// Validate configuration data using the generic validator
	validator := validators.NewModelValidator[models.Configuration]()
	config, err := validator.ValidateFromMap(configData)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	// Set account_id from context if not provided in payload
	if config.AccountID == 0 {
		config.AccountID = utils.GetCurrentAccountID(ctx)
	}

	// Create configuration using the service
	createdConfig, err := h.configurationService.CreateConfiguration(ctx, configData)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	utils.WriteJSONResponse(w, http.StatusCreated, createdConfig, "Configuration created successfully")
}

// func (h *ConfigurationHandler) UpdateConfigurationHandler(w http.ResponseWriter, r *http.Request) {
// 	payload, err := utils.ParseJSONBody(r)
// 	if err != nil {
// 		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	configData, exists := payload["configuration"].(map[string]any)
// 	if !exists || len(configData) == 0 {
// 		utils.WriteErrorResponse(w, http.StatusBadRequest, "Configuration data can't be blank")
// 		return
// 	}

// 	ctx := r.Context()
// 	config, err := h.configurationService.UpdateConfiguration(ctx, configData)
// 	if err != nil {
// 		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}
// 	utils.WriteJSONResponse(w, http.StatusOK, config, "Configuration updated successfully")
// }
