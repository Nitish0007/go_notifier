package middlewares

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
)

func AuthenticateRequest(conn *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized request")
				return
			}

			authKey := r.Header.Get("Authorization")
			accountIDStr := chi.URLParam(r, "account_id")
			accountID, err := strconv.Atoi(accountIDStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid account ID")
				return
			}

			ctx := r.Context()
			repo := repositories.NewApiKeyRepository(conn)

			apiKey, err := repo.FindByAccountID(ctx, accountID)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			if apiKey.Key == "" || apiKey.Key != authKey {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Api Key used")
			}

			ctx = utils.SetCurrentAccountID(ctx, accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
