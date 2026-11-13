package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error retrieving chirps"))
		return
	}
	responseChirps := []Chirp{}
	for _, c := range chirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			UserID:    c.UserID,
			Body:      c.Body,
		})
	}

	dat, err := json.Marshal(responseChirps)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling response Chirps"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	parsedUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid UUID"))
		return
	}
	c, err := cfg.db.GetChirp(r.Context(), parsedUUID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(404)
		w.Write([]byte("chirp doesn't exist"))
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error retrieving chirp"))
		return
	}
	responseChirp := Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		UserID:    c.UserID,
		Body:      c.Body,
	}
	dat, err := json.Marshal(responseChirp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling reponse chirp"))
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
