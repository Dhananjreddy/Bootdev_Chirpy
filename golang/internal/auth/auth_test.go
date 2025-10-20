package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/Dhananjreddy/Bootdev_Chirpy/golang/internal/auth"
)


func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecure123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	ok, err := auth.CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !ok {
		t.Error("Expected password match, got mismatch")
	}

	ok, _ = auth.CheckPasswordHash("wrongpassword", hash)
	if ok {
		t.Error("Expected mismatch for wrong password, got match")
	}
}


func TestMakeAndValidateJWT(t *testing.T) {
	tokenSecret := "mysecretkey"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedID, err := auth.ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	tokenSecret := "secret1"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Error("Expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	tokenSecret := "secret123"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, tokenSecret, -time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headerVal string
		wantToken string
		wantErr   bool
	}{
		{"valid", "Bearer abc123", "abc123", false},
		{"missing header", "", "", true},
		{"wrong scheme", "Token abc123", "", true},
		{"malformed header", "Bearer", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.headerVal != "" {
				h.Set("Authorization", tt.headerVal)
			}

			got, err := auth.GetBearerToken(h)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
			if got != tt.wantToken {
				t.Errorf("expected token=%q, got %q", tt.wantToken, got)
			}
		})
	}
}
