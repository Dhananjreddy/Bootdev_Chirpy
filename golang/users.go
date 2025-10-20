package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/database"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/auth"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	HashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password", err)
    	return
	}

	user, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: HashedPassword,
		Email: params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (apiCfg *apiConfig) handlerUpdatePassword(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type response struct {
		User
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
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
    }

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash password", err)
		return
    }
	
	updateParams := database.UpdatePasswordParams{
		HashedPassword: hashedPassword,
		Email: params.Email,
		ID: UserID,
	}

	user, err := apiCfg.db.UpdatePassword(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 500, "Couldn't update password", err)
		return
	}

	respondWithJSON(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}