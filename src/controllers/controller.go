package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
	"keycloak-go-backend/src/services"
)

type doc struct {
	Id   string    `json:"id"`
	Num  string    `json:"num"`
	Date time.Time `json:"date"`
}

type Controller struct {
	keycloak *services.Keycloak
}

type loginRequest struct {
	Username string
	Password string
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func NewController(keycloak *services.Keycloak) *Controller {
	return &Controller{
		keycloak: keycloak,
	}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	rq := &loginRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var decoder fastjson.Parser

	value, err := decoder.ParseBytes(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	rq.Username = string(value.GetStringBytes("username"))
	rq.Password = string(value.GetStringBytes("password"))

	jwt, err := c.keycloak.Gocloak.Login(context.Background(),
		c.keycloak.ClientId,
		c.keycloak.ClientSecret,
		c.keycloak.Realm,
		rq.Username,
		rq.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rs := &loginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

	rsJs, err := json.Marshal(rs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rsJs)
}

func (c *Controller) GetDocs(w http.ResponseWriter, r *http.Request) {
	rs := []*doc{
		{
			Id:   "1",
			Num:  "ABC-123",
			Date: time.Now().UTC(),
		},
		{
			Id:   "2",
			Num:  "ABC-456",
			Date: time.Now().UTC(),
		},
	}

	rsJs, _ := json.Marshal(rs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rsJs)

}
