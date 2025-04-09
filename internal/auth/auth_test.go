package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCompareHashPassord(t *testing.T) {
	password := "12345678"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error creating hash from password: %s", err)
	}

	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("password didn't produce expected hash: %s", err)
	}

	err = CheckPasswordHash(hash, "12345")
	if err == nil {
		t.Errorf("expected an error since the hash doesn't match the password")
	}
}

func TestValidJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	expiresIn := 10 * time.Second
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Errorf("error signing new token: %s", err)
	}

	id, err := ValidateJWT(token, secret)
	if err != nil || id != userID {
		t.Errorf("token should be valid: %s", err)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	expiresIn := -10 * time.Second
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Errorf("error signing new token: %s", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expired token should not be valid")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	expiresIn := 10 * time.Second
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Errorf("error signing new token: %s", err)
	}

	_, err = ValidateJWT(token, "wrong secret")
	if err == nil {
		t.Errorf("should not be able to validate the token with a wrong secret")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := make(http.Header)
	token := "some-token"
	header.Set("Authorization", "Bearer " + token)

	parsedToken, err := GetBearerToken(header)
	if err != nil || token != parsedToken {
		t.Errorf("error parsing bearer token")
	}
}