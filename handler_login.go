package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-http-servers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	loginCreds := Body{}
	err := json.NewDecoder(r.Body).Decode(&loginCreds)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("eror decoding login"))
		return
	}
	user, err := cfg.db.GetUserEmail(r.Context(), loginCreds.Email)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(401)
		w.Write([]byte("Email does not exist in database"))
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error retrieving user"))
		return
	}
	match, err := auth.CheckPasswordHash(loginCreds.Password, user.Password)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error checking password"))
		return
	}
	if !match {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
		return
	}
	responseUser := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	dat, err := json.Marshal(responseUser)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error marshalling data"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
