package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Claims struct {
		Data []byte `json:"data,omitempty"`
		jwt.RegisteredClaims
	}

	TokenManager[T any] interface {
		GenerateToken(payload T) (string, error)
		ValidateToken(tokenString string) (T, error)
	}

	Manager[T any] struct {
		secret   []byte
		tokenTTL time.Duration
	}
)

func NewManager[T any](secret string, tokenTTL time.Duration) *Manager[T] {
	return &Manager[T]{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

// GenerateToken generates a JWT token with a given payload.
func (m Manager[T]) GenerateToken(payload T) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token payload: %w", err)
	}

	expiration := jwt.NewNumericDate(time.Now().Add(m.tokenTTL))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Data:             data,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expiration},
	})

	return token.SignedString(m.secret)
}

// ValidateToken validates the given JWT token and returns a payload.
func (m Manager[T]) ValidateToken(tokenString string) (T, error) {
	var payload T

	claims := new(Claims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.secret, nil
	})
	if err != nil {
		return payload, err
	}

	if !token.Valid {
		return payload, errors.New("invalid token")
	}

	if err = json.Unmarshal(claims.Data, &payload); err != nil {
		return payload, fmt.Errorf("failed to unmarshal token payload: %w", err)
	}

	return payload, nil
}
