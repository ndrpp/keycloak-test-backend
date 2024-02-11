package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"keycloak-go-backend/src/services"
	"keycloak-go-backend/src/utils"
	"net/http"
	"strings"
)

func extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func VerifyToken(logger *utils.Logger, next http.Handler) http.Handler {
	Keycloak := *services.NewKeycloak()

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			token = extractBearerToken(token)

			if token == "" {
				http.Error(w, "Bearer Token missing", http.StatusUnauthorized)
				return
			}

			result, err := Keycloak.Gocloak.RetrospectToken(context.Background(), token, Keycloak.ClientId, Keycloak.ClientSecret, Keycloak.Realm)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			jwt, _, err := Keycloak.Gocloak.DecodeAccessToken(context.Background(), token, Keycloak.Realm)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			_, err = json.Marshal(jwt)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to marshal jwt response: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			if !*result.Active {
				http.Error(w, "Invalid or expired Token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
}

func CorsMiddeware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
