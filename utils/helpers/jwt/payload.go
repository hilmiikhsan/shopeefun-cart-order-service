package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	Email string
	jwt.RegisteredClaims
}

func NewPayload(email string, duration time.Duration) (*Payload, error) {
	usrEmail, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error on generating uuid")
		return nil, err
	}

	timeNow := time.Now()
	payload := &Payload{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
			Issuer:    "user_login",
			Subject:   "shopifun",
			ID:        usrEmail.String(),
		},
	}
	return payload, nil
}
