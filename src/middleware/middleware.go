package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"keycloak-go-backend/src/services"
	"net/http"
	"strings"
)

type KeyCloakMiddleware struct {
	Keycloak *services.Keycloak
}

func NewMiddleware(keycloak *services.Keycloak) *KeyCloakMiddleware {
	return &KeyCloakMiddleware{Keycloak: keycloak}
}

func (auth *KeyCloakMiddleware) extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func (auth *KeyCloakMiddleware) VerifyToken(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token = auth.extractBearerToken(token)

		if token == "" {
			http.Error(w, "Bearer Token missing", http.StatusUnauthorized)
			return
		}

		result, err := auth.Keycloak.Gocloak.RetrospectToken(context.Background(), token, auth.Keycloak.ClientId, auth.Keycloak.ClientSecret, auth.Keycloak.Realm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		jwt, _, err := auth.Keycloak.Gocloak.DecodeAccessToken(context.Background(), token, auth.Keycloak.Realm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		jwtj, err := json.Marshal(jwt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal jwt response: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		fmt.Printf("token: %v\n", string(jwtj))

		if !*result.Active {
			http.Error(w, "Invalid or expired Token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
