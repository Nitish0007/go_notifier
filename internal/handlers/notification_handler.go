package handlers

import (
	"log"
	"net/http"

	"github.com/Nitish0007/go_notifier/internal/services"
	"github.com/Nitish0007/go_notifier/utils"
)

type NotificationHandler struct{
	notificationService *services.NotificationService
}

func NewNotificationHandler(s *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: s,
	}
}

func (h *NotificationHandler) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	
	payload, err := utils.ParseJSONBody(r)

	if err != nil {
		log.Printf(">>>> ERROR: %v", err.Error())
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	notificationData, exists := payload["notification"].(map[string]any)
	if !exists {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Notification data can't be blank")
		return
	}

	ctx := r.Context()
	n, err := h.notificationService.CreateNotification(ctx, notificationData)
	if err != nil {
		log.Printf(">>>> ERROR: %v", err.Error())
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.notificationService.SendOrScheduleNotification(ctx, n)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil ,"Notification enqueued successfully")
}

func (h *NotificationHandler) SendBulkNotificationHandler(w http.ResponseWriter, r *http.Request) {

} 