package main

import (
	"errors"
	"net/http"

	"github.com/bootdotdev/learn-http-servers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "couldn't parse path value", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "couldn't receive token from bearer", err)
		return
	}
	if tokenString == "" {
		respondWithError(w, 401, "no auth header", err)
		return
	}
	id, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "couldn't validate", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, "couldn't find chirp", err)
		return
	}
	if chirp.UserID != id {
		respondWithError(w, 403, "user not owner of chirp", errors.New("unauthorized"))
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 500, "couldn't delete chirp", err)
		return
	}
	w.WriteHeader(204)
}
