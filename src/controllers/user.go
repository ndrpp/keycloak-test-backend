package controllers

import (
	"context"
	"net/http"
	"time"

	"keycloak-go-backend/src/services"
	"keycloak-go-backend/src/utils"
)

type doc struct {
	Id   string    `json:"id"`
	Num  string    `json:"num"`
	Date time.Time `json:"date"`
}

type UserController struct {
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

func NewUserController(keycloak *services.Keycloak) *UserController {
	return &UserController{
		keycloak: keycloak,
	}
}

func (c *UserController) Login(logger *utils.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rq := &loginRequest{}
			rq, err := utils.Decode[*loginRequest](r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

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
			err = utils.Encode[*loginResponse](w, r, http.StatusOK, rs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func (c *UserController) GetDocs(logger *utils.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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

			err := utils.Encode[[]*doc](w, r, http.StatusOK, rs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func (c *UserController) HandleHealthZ(logger *utils.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := utils.Encode[string](w, r, http.StatusOK, "Alive and well.")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}
