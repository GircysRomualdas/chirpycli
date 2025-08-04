package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Couldn't hash password: %v", err)
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("Password does not match")
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signature, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Couldn't sign JWT: %v", err)
	}
	return signature, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claim := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't parse JWT: %v", err)
	}
	if !token.Valid {
		return uuid.Nil, errors.New("JWT is invalid")
	}
	userID, err := uuid.Parse(claim.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't parse user ID from JWT: %v", err)
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawAuthToken := headers.Get("Authorization")
	if rawAuthToken == "" {
		return "", errors.New("Couldn't find authorization token")
	}
	if !strings.HasPrefix(rawAuthToken, "Bearer ") {
		return "", errors.New("Authorization token must start with 'Bearer '")
	}
	authToken := strings.TrimSpace(strings.TrimPrefix(rawAuthToken, "Bearer "))
	return authToken, nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Couldn't generate random bytes: %v", err)
	}
	return hex.EncodeToString(b), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	rawAPIKey := headers.Get("Authorization")
	if rawAPIKey == "" {
		return "", errors.New("Couldn't find API key")
	}
	if !strings.HasPrefix(rawAPIKey, "ApiKey ") {
		return "", errors.New("API key must start with 'ApiKey '")
	}
	apiKey := strings.TrimSpace(strings.TrimPrefix(rawAPIKey, "ApiKey "))
	return apiKey, nil
}
