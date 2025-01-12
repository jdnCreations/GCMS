package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "gcms",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// pass pointer so modifications will be made to struct
	claims := &jwt.RegisteredClaims{} 

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error)  {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// get auth from headers
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no Bearer")
	}

	if !strings.HasPrefix(auth, "Bearer") {
		return "", errors.New("invalid Bearer format")
	}

	// get the data after "Bearer "
	trimmed := strings.TrimPrefix(auth, "Bearer ")
	if trimmed == "" {
		return "", errors.New("no Bearer")
	}
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return "", errors.New("invalid bearer")
	}
	
	return trimmed, nil
}