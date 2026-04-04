package middlewares

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
)

func AuthenticateRequest(conn *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				api.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized request")
				return
			}

			authKey := r.Header.Get("Authorization")
			accountIDStr := chi.URLParam(r, "account_id")
			accountID, err := strconv.Atoi(accountIDStr)
			if err != nil {
				api.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid account ID")
				return
			}

			ctx := r.Context()
			repo := apiKey.NewApiKeyRepository(conn)

			apiKey, err := repo.FindByAccountID(ctx, accountID)
			if err != nil {
				api.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			if apiKey.Key == "" || apiKey.Key != authKey {
				api.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Api Key used")
				return
			}

			ctx = sharedhelper.SetCurrentAccountID(ctx, accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
