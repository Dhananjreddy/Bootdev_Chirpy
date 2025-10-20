package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"time"
	"errors"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/database"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/auth"
	"database/sql"
	"sort"
)


type newChirp struct {
	Id 			uuid.UUID			`json:"id"`
	CreatedAt   time.Time 			`json:"created_at"`
	UpdatedAt   time.Time 			`json:"updated_at"`
	Body 		string 				`json:"body"`
	UserID		uuid.UUID			`json:"user_id"`
}

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){
	
	type chirpParams struct{
		Body   string       `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	UserID, err := auth.ValidateJWT(token, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}


	decoder := json.NewDecoder(r.Body)
    params := chirpParams{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
    }

	params.Body, err = validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, "Invalid Chirp", err)
	}

	chirp, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Error Creating Chirp", err)
	}

	respondWithJSON(w, 201, newChirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})

	return
}

func (apiCfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request){

	s := r.URL.Query().Get("author_id")
	if s != "" {
		id, err := uuid.Parse(s)
    	if err != nil {
			respondWithError(w, 500, "Couldn't parse author_id", err)
			return
    	}

		sort := r.URL.Query().Get("sort")
		if sort == "desc"{
			chirps, err := apiCfg.db.GetChirpsByUserIDDesc(r.Context(), id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					respondWithError(w, 404, "No chirps found", err)
					return
				}
				respondWithError(w, 500, "Error fetching chirps from database", err)
				return
			}
			var chirpsOut []newChirp

			for _, chirp := range chirps {
				chirpsOut = append(chirpsOut, newChirp{
				Id: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
				})
			}

			respondWithJSON(w, 200, chirpsOut)
			return
		}

		chirps, err := apiCfg.db.GetChirpsByUserIDAsc(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, 404, "No chirps found", err)
				return
			}
			respondWithError(w, 500, "Error fetching chirps from database", err)
			return
		}
		var chirpsOut []newChirp

		for _, chirp := range chirps {
			chirpsOut = append(chirpsOut, newChirp{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
			})
		}

		respondWithJSON(w, 200, chirpsOut)
		return
	}

	chirps, err := apiCfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error fetching chirps from database", err)
		return
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc"{
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
        	return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
    	})
	}

	var chirpsOut []newChirp

	for _, chirp := range chirps {
		chirpsOut = append(chirpsOut, newChirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, 200, chirpsOut)
}

func (apiCfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request){

	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr); if err != nil {
		respondWithError(w, 400, "Invalid Chirp id", nil)
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	respondWithError(w, 404, "Chirp not found", err)
        	return
    	}
    	respondWithError(w, 500, "Error fetching chirp from database", err)
    	return
	}

	chirpOut := newChirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 200, chirpOut)
	return
}

func (apiCfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request){
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr); if err != nil {
		respondWithError(w, 400, "Invalid Chirp id", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	UserID, err := auth.ValidateJWT(token, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	respondWithError(w, 404, "Chirp not found", err)
        	return
    	}
    	respondWithError(w, 500, "Error fetching chirp from database", err)
    	return
	}

	if chirp.UserID != UserID {
		respondWithError(w, 403, "Not chirp Author, Cannot delete", err)
		return
	}
	
	err = apiCfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	respondWithError(w, 404, "Chirp not found", err)
        	return
    	}
    	respondWithError(w, 500, "Error deleting chirp from database", err)
    	return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func validateChirp(chirp string) (string, error){
    
	const maxChirpLength = 140
	if len(chirp) > maxChirpLength {
		return "", errors.New("Chirp of invalid Length")
	}

	return checkProfanity(chirp), nil
}

func checkProfanity(s string) string {
	words := strings.Split(s, " ")
	for i, word := range(words){
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	final := strings.Join(words, " ")
	return final
}