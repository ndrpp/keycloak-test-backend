package services

import "github.com/Nerzal/gocloak/v13"

type Keycloak struct {
	Gocloak      *gocloak.GoCloak
	ClientId     string
	ClientSecret string
	Realm        string
}

func NewKeycloak() *Keycloak {
	return &Keycloak{
		Gocloak:      gocloak.NewClient("http://127.0.0.1:8080"),
		ClientId:     "my-go-service",
		ClientSecret: "ZapjG4VaQf4bEBwCrFP0zhMVehP3jbh6",
		Realm:        "go-test",
	}
}
