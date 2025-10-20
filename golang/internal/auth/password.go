package auth

import (
	"github.com/alexedwards/argon2id"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"errors"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error parsing hash: %v", err)
		
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error){
	return argon2id.ComparePasswordAndHash(password, hash)
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chirpy",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error making JWT: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token: %v", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header found")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}


func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}


// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included in request")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}