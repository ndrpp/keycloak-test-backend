package services

import (
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type Keycloak struct {
	Gocloak      *gocloak.GoCloak
	ClientId     string
	ClientSecret string
	Realm        string
}

func NewKeycloak() *Keycloak {
	return &Keycloak{
		Gocloak:      gocloak.NewClient("http://127.0.0.1:8080"),
		ClientId:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
	}
}
