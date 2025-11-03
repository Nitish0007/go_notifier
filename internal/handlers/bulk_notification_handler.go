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

	notificationsRaw, exists := payload["notifications"]
	if !exists {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Notifications data is required")
		return
	}

	notifications, err := json.Marshal(notificationsRaw)  // converting to json format
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Unable to marshal notifications")
		return
	}

	var notificationsList []map[string]any
	err = json.Unmarshal(notifications, &notificationsList)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Unable to unmarshal notifications")
		return
	}

	response, err := h.bulkNotificationService.CreateBulkNotifications(ctx, notificationsList) 
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response, "")
}