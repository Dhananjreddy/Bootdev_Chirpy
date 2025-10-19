package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"time"
	"errors"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/database"

)

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){
	type chirpParams struct{
		Body   string       `json:"body"`
		UserID uuid.UUID	`json:"user_id"`
	}

	type newChirp struct {
		Id 			uuid.UUID			`json:"id"`
		CreatedAt   time.Time 			`json:"created_at"`
		UpdatedAt   time.Time 			`json:"updated_at"`
		Body 		string 				`json:"body"`
		UserID		uuid.UUID			`json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
    params := chirpParams{}
    err := decoder.Decode(&params)
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
		UserID: params.UserID,
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