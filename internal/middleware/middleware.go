package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

const API_KEY_HEADER = "api_key"
const EXPECTED_API_KEY = "apitest"

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(API_KEY_HEADER)
		if apiKey == "" {
			w.WriteHeader(401)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(model.Response{
				Code:    401,
				Type:    "Unauthorised",
				Message: "Missing API key",
			})
			return
		}
		if apiKey != EXPECTED_API_KEY {
			w.WriteHeader(403)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(model.Response{
				Code:    403,
				Type:    "Forbidden",
				Message: "Wrong API key",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
