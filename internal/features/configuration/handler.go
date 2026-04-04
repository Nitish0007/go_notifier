package configuration

import (
	// "log"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
	"github.com/Nitish0007/go_notifier/internal/shared/validators"
)

type ConfigurationHandler struct {
	configurationService *ConfigurationService
}

func NewConfigurationHandler(s *ConfigurationService) *ConfigurationHandler {
	return &ConfigurationHandler{
		configurationService: s,
	}
}

// for getting configurations
func (h *ConfigurationHandler) GetConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	configs, err := h.configurationService.GetConfigurations(ctx, accID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusOK, configs, "Configurations fetched successfully")
}

// for creating configuration
func (h *ConfigurationHandler) CreateConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqData, err := api.ParseAndValidateRequest[CreateConfigurationRequest](r)
	if err != nil {
		if errors.Is(err, &validators.ValidationError{Message: "AccountID is required", Field: "account_id", Tag: "required"}) {
			reqData.Configuration.AccountID = sharedhelper.GetCurrentAccountID(ctx)
		} else {
			api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		return
	}

	// Validate flat configuration_data by config_type (smtp / in_app / future types)
	if reqData.Configuration.ConfigurationData != nil {
		validated, err := ValidateConfigurationDataByType(reqData.Configuration.ConfigType, reqData.Configuration.ConfigurationData)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		reqData.Configuration.ConfigurationData = validated
	}

	// create configuration using the service
	createdConfig, err := h.configurationService.CreateConfiguration(ctx, reqData)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusCreated, createdConfig, "Configuration created successfully")

	
}

func (h *ConfigurationHandler) DeleteConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	confID, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	if accID == -1 {
		api.WriteErrorResponse(w, http.StatusUnauthorized, "Unknown account ID")
		return
	}

	cidInt, err := strconv.Atoi(confID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	err = h.configurationService.DeleteConfiguration(ctx, accID, cidInt)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	api.WriteResponse(w, http.StatusOK, nil, "Configuration deleted successfully")
}

func (h *ConfigurationHandler) UpdateConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	confID, err := api.GetPathParam(r, "id")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cidInt, err := strconv.Atoi(confID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	reqData, err := api.ParseAndValidateRequest[UpdateConfigurationRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	reqData.Configuration.ID = cidInt
	reqData.Configuration.AccountID = sharedhelper.GetCurrentAccountID(r.Context())

	// Validate flat configuration_data by config_type when present
	if reqData.Configuration.ConfigurationData != nil {
		validated, err := ValidateConfigurationDataByType(reqData.Configuration.ConfigType, reqData.Configuration.ConfigurationData)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		reqData.Configuration.ConfigurationData = validated
	}

	ctx := r.Context()
	updatedConfig, err := h.configurationService.UpdateConfiguration(ctx, reqData)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusOK, updatedConfig, "Configuration updated successfully")
}
