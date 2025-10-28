package handlers

import (
	"net/http"

	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/services"
)

type BulkNotificationHandler struct {
	bulkNotificationService *services.BulkNotificationService
}

func NewBulkNotificationHandler(s *services.BulkNotificationService) *BulkNotificationHandler {
	return &BulkNotificationHandler{
		bulkNotificationService: s,
	}
}

func (h *BulkNotificationHandler) CreateBulkNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.ParseJSONBody(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Unable to parse payload")
		return
	}

	ctx := r.Context()
	response, err := h.bulkNotificationService.CreateBulkNotifications(ctx, payload["notifications"].([]map[string]any)) 
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, "invalid payloads")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response, "")
}