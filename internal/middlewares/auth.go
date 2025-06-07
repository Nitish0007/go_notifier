package middlewares

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	
	"github.com/Nitish0007/go_notifier/utils"
	"github.com/Nitish0007/go_notifier/internal/repositories"
)

func AuthenticateRequest(conn *pgx.Conn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			if r.Header.Get("Authorization") == "" {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized request")
				return
			}

			authKey := r.Header.Get("Authorization")
			accountIDStr := chi.URLParam(r, "account_id")
			accountID, err := strconv.Atoi(accountIDStr)

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

			next.ServeHTTP(w, r)
		})
	}
}
	