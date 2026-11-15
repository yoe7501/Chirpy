package main

import (
	"net/http"
	"time"

	"github.com/bootdotdev/learn-http-servers/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}
	oldToken, err := cfg.db.GetToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 404, "Couln't find refresh token", err)
		return
	}

	if oldToken.RevokedAt.Valid {
		respondWithError(w, 401, "Token revoked", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}
