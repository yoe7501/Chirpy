package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-http-servers/internal/auth"
	"github.com/bootdotdev/learn-http-servers/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateCred(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error finding token", err)
	}
	if tokenString == "" {
		respondWithError(w, 401, "Auth header missing", err)
	}
	userId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error validating jwt", err)
	}
	var newBody body
	err = json.NewDecoder(r.Body).Decode(&newBody)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body", err)
	}
	hashed, err := auth.HashPassword(newBody.Password)
	if err != nil {
		respondWithError(w, 500, "coudln't hash password", err)
	}

	updateParams := database.UpdateUserParams{
		Email:    newBody.Email,
		Password: hashed,
		ID:       userId,
	}
	user, err := cfg.db.UpdateUser(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 500, "couldn't update user", err)
	}

	responseUser := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	dat, err := json.Marshal(responseUser)
	if err != nil {
		respondWithError(w, 500, "couldn't marshal responseUser", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
