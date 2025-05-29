package middlewares

import (
	"net/http"

	"github.com/Nitish0007/go_notifier/utils"
)

func AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Header.Get("Authorization") == "" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized request")
			return
		}

		next.ServeHTTP(w, r)
	})
}