package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is valid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	fmt.Println("uuid ", tokenID)
	payload := &Payload{

		username,
		jwt.RegisteredClaims{
			ID:        tokenID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload, nil
}

// force function Valid

func (payload *Payload) Valid() error {
	if time.Now().After(payload.RegisteredClaims.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
