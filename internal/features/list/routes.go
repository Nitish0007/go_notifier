package list

import (
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
)

func RegisterListRoutes(conn *gorm.DB, r chi.Router, h *ListHandler) {
	r.Get("/lists", h.GetListsHandler)
	r.Post("/lists", h.CreateListHandler)
	r.Post("lists/{id}/subscribe", h.SubscribeToListHandler)
	r.Post("/lists/{id}/manage_subscription", h.ManageListSubscriptionHandler)
	// r.Get("/lists/{id}", h.GetListByIdHandler)
	// r.Get("/lists/{uuid}", h.GetListByUUIDHandler)
}