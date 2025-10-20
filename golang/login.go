package main

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/auth"
	"database/sql"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type loginParams struct{
		Password string 	`json:"password"`
		Email 	 string 	`json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	params := loginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
    }
	
	user, err := apiCfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
    	if errors.Is(err, sql.ErrNoRows) {
        	respondWithError(w, 404, "User not found", err)
        	return
    	}
    	respondWithError(w, 500, "Error fetching User from database", err)
    	return
	}

	
	login, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 500, "Error verifying password", err)
    	return
	}

	if !login {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, 200, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
	return
}