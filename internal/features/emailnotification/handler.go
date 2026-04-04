package emailnotification

import (
	// "log"
	"net/http"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

type EmailNotificationHandler struct {
	notificationService *EmailNotificationService
}

func NewEmailNotificationHandler(s *EmailNotificationService) *EmailNotificationHandler {
	return &EmailNotificationHandler{
		notificationService: s,
	}
}

// func (h *EmailNotificationHandler) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {

// 	payload, err := api.ParseJSONBody(r)

// 	if err != nil {
// 		log.Printf(">>>> ERROR: %v", err.Error())
// 		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	notificationData, exists := payload["notification"].(map[string]any)
// 	if !exists {
// 		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Notification data can't be blank")
// 		return
// 	}

// 	ctx := r.Context()
// 	n, err := h.notificationService.CreateNotification(ctx, notificationData)
// 	if err != nil {
// 		log.Printf(">>>> ERROR: %v", err.Error())
// 		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}

// 	err = h.notificationService.SendOrScheduleNotification(ctx, n)
// 	if err != nil {
// 		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}

// 	api.WriteResponse(w, http.StatusOK, nil, "Notification enqueued successfully")
// }

func (h *EmailNotificationHandler) CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[CreateEmailCampaignRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	email_notif, err := h.notificationService.CreateEmailCampaign(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, email_notif, "Email Notification Campaign created successfully")
}

func (h *EmailNotificationHandler) CreateEmailTransactionalHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseAndValidateRequest[CreateEmailTransactionalRequest](r)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	trans_notif, err := h.notificationService.CreateEmailTransactionalCampaign(ctx, payload)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusCreated, trans_notif, "Email Transactional Notification created successfully")
}

// TODO: add pagination and filtering
// index api for notifications in context of account
func (h *EmailNotificationHandler) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := sharedhelper.GetCurrentAccountID(ctx)
	list, err := h.notificationService.GetNotifications(ctx, accID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	api.WriteResponse(w, http.StatusOK, list, "Success")
}

// func (h *EmailNotificationHandler) SendNotificationByIDHandler(w http.ResponseWriter, r *http.Request) {
// 	payload, err := api.ParseJSONBody(r)
// 	if err != nil {
// 		api.WriteErrorResponse(w, http.StatusBadRequest, "Unable to parse payload")
// 		return
// 	}

// 	nID, exists := payload["notification_id"].(string)
// 	if !exists || !sharedhelper.IsValidUUID(nID) {
// 		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid notification id")
// 		return
// 	}

// 	ctx := r.Context()
// 	accID := sharedhelper.GetCurrentAccountID(ctx)

// 	n, err := h.notificationService.GetNotificationById(ctx, nID, accID)
// 	if err != nil {
// 		log.Printf("ERROR!: %v", err)
// 		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, "not able to fetch notification")
// 		return
// 	}

// 	err = h.notificationService.SendOrScheduleNotification(ctx, n)
// 	if err != nil {
// 		log.Printf("ERROR!: %v", err)
// 		api.WriteErrorResponse(w, http.StatusUnprocessableEntity, "not able to schedule notification")
// 		return
// 	}

// 	api.WriteResponse(w, http.StatusOK, nil, "Notification enqued successfully if its send_at time is in past or 10 within 10 minutes in future")
// }