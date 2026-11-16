package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-http-servers/internal/auth"
	"github.com/bootdotdev/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID         uuid.UUID `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Token      string    `json:"token"`
		Refresh    string    `json:"refresh_token"`
		IsChirpRed bool      `json:"is_chirpy_red"`
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
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(3600))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error making JWT token"))
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error making refresh token"))
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}
	responseUser := response{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Email:      user.Email,
		Token:      token,
		Refresh:    refreshToken,
		IsChirpRed: user.IsChirpyRed,
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
