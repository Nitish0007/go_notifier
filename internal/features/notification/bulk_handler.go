package notification

// import (
// 	"net/http"

// 	"github.com/Nitish0007/go_notifier/utils"
// )

// type BulkNotificationHandler struct {
// 	bulkNotificationService *NotificationService
// }

// func NewBulkNotificationHandler(s *NotificationService) *BulkNotificationHandler {
// 	return &BulkNotificationHandler{
// 		bulkNotificationService: s,
// 	}
// }

// func (h *BulkNotificationHandler) CreateBulkNotificationsHandler(w http.ResponseWriter, r *http.Request) {
// 	payload, err := utils.ParseJSONBody(r)
// 	if err != nil {
// 		utils.WriteErrorResponse(w, http.StatusBadRequest, "Unable to parse payload")
// 		return
// 	}

// 	ctx := r.Context()
// 	response, err := h.bulkNotificationService.CreateBulkNotifications(ctx, payload) 
// 	if err != nil {
// 		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}

// 	utils.WriteJSONResponse(w, http.StatusOK, response, "")
// }