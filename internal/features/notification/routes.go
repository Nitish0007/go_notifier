package notification

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterNotificationRoutes(conn *gorm.DB, r chi.Router, h *NotificationHandler) {
	// protected routes
	r.Post("/notify", h.SendNotificationHandler) // create and send notification and return the notification id
	// r.Post("/bulk_notify", h.SendBulkNotificationHandler)
	r.Post("/notification", h.CreateNotificationHandler) // create notification and return its object
	r.Post("/trigger", h.SendNotificationByIDHandler)     // send notification by id
	r.Get("/notifications", h.GetNotificationsHandler)    // get notifications in context of account
}

// func RegisterBulkNotificationRoutes(conn *gorm.DB, r chi.Router, h *BulkNotificationHandler) {
// 	r.Post("/bulk_notifications", h.CreateBulkNotificationsHandler) // create bulk notifications and return list of errors if any or return successfully created
// }
