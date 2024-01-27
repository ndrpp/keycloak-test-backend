package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type keyCloakMiddleware struct {
	keycloak *keycloak
}

func newMiddleware(keycloak *keycloak) *keyCloakMiddleware {
	return &keyCloakMiddleware{keycloak: keycloak}
}

func (auth *keyCloakMiddleware) extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func (auth *keyCloakMiddleware) verifyToken(next http.Handler) http.Handler {

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
		fmt.Println("extracted token: ", token)

		result, err := auth.keycloak.gocloak.RetrospectToken(context.Background(), token, auth.keycloak.clientId, auth.keycloak.clientSecret, auth.keycloak.realm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}
		fmt.Println("Result from retrospect: ", result)

		jwt, _, err := auth.keycloak.gocloak.DecodeAccessToken(context.Background(), token, auth.keycloak.realm)
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
