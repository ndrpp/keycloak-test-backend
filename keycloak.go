package main

import "github.com/Nerzal/gocloak/v7"

type keycloak struct {
	gocloak      gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func newKeycloak() *keycloak {
	return &keycloak{
		gocloak:      gocloak.NewClient("http://localhost:8086"),
		clientId:     "my-go-service",
		clientSecret: "0EP350HKcb8zDneG2gsSb5KwCtRBphUK",
		realm:        "go-test",
	}
}
