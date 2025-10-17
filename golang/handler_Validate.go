package main

import(
	"net/http"
	"encoding/json"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

	type ReturnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
    }

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}

	returnVal := ReturnVal{
		CleanedBody: checkProfanity(params.Body),
	}
	respondWithJSON(w, 200, returnVal)
	return
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