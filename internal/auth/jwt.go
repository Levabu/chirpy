package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TOKEN_EXP_TIME = time.Hour * 1
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(
		tokenString, 
		&jwt.RegisteredClaims{}, 
		func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok || !t.Valid {
		return uuid.Nil, jwt.ErrSignatureInvalid
	}

	return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	token, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok || token == "" {
		return "", errors.New("no authorization or bearer token was provided")
	}
	return token, nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	apiKey, ok := strings.CutPrefix(authHeader, "ApiKey ")
	if !ok || apiKey == "" {
		return "", errors.New("no api key was provided")
	}
	return apiKey, nil
}

