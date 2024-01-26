package main

import (
	"context"
	"encoding/json"
	//"fmt"
	//"io"
	"net/http"
	"time"
	//"github.com/valyala/fastjson"
)

type doc struct {
	Id   string    `json:"id"`
	Num  string    `json:"num"`
	Date time.Time `json:"date"`
}

type controller struct {
	keycloak *keycloak
}

type loginRequest struct {
	Username string
	Password string
}

type loginResponse struct {
	accessToken  string
	refreshToken string
	expiresIn    int
}

func newController(keycloak *keycloak) *controller {
	return &controller{
		keycloak: keycloak,
	}
}

func (c *controller) login(w http.ResponseWriter, r *http.Request) {
	//rq := &loginRequest{}

	//body, err := io.ReadAll(r.Body)
	//fmt.Println("Received body on login...")
	//if err != nil {
	//	fmt.Println("Error reading request body: ", err)
	//	http.Error(w, "Bad Request", http.StatusBadRequest)
	//	return
	//}
	//var decoder fastjson.Parser

	//value, err := decoder.ParseBytes(body)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//}

	//rq.username = string(value.GetStringBytes("username"))
	//rq.password = string(value.GetStringBytes("password"))
	rq := &loginRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := c.keycloak.gocloak.Login(context.Background(),
		c.keycloak.clientId,
		c.keycloak.clientSecret,
		c.keycloak.realm,
		rq.Username,
		rq.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rs := &loginResponse{
		accessToken:  jwt.AccessToken,
		refreshToken: jwt.RefreshToken,
		expiresIn:    jwt.ExpiresIn,
	}

	rsJs, _ := json.Marshal(rs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rsJs)
}

func (c *controller) getDocs(w http.ResponseWriter, r *http.Request) {
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
