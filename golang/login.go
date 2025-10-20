package main

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/auth"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/database"
	"database/sql"
	"time"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type loginParams struct{
		Password 				string 					`json:"password"`
		Email 	 				string 					`json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(
		user.ID,
		apiCfg.secret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token: accessToken,
		RefreshToken: refreshToken,
	})
}