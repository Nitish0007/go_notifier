package emailnotification

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterEmailNotificationRoutes(conn *gorm.DB, r chi.Router, h *EmailNotificationHandler) {
	// protected routes
	r.Post("/email_notifications/transactional", h.CreateEmailTransactionalHandler)
	r.Post("/email_notifications/campaign", h.CreateNotificationHandler)
	// r.Post("/notify", h.SendNotificationHandler) // create and send notification and return the notification id
	// r.Post("/bulk_notify", h.SendBulkNotificationHandler)
	// r.Post("/trigger", h.SendNotificationByIDHandler)     // send notification by id
	// r.Get("/notifications", h.GetNotificationsHandler)    // get notifications in context of account
}

// func RegisterBulkNotificationRoutes(conn *gorm.DB, r chi.Router, h *BulkNotificationHandler) {
// 	r.Post("/bulk_notifications", h.CreateBulkNotificationsHandler) // create bulk notifications and return list of errors if any or return successfully created
// }