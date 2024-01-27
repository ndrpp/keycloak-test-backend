package main

import "github.com/Nerzal/gocloak/v13"

type keycloak struct {
	gocloak      *gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func newKeycloak() *keycloak {
	return &keycloak{
		gocloak:      gocloak.NewClient("http://127.0.0.1:8080"),
		clientId:     "my-go-service",
		clientSecret: "ZapjG4VaQf4bEBwCrFP0zhMVehP3jbh6",
		realm:        "go-test",
	}
}
