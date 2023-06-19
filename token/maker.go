package token

import (
	"time"
)

//Maker interface manage token

type Maker interface {
	// CreateToken from username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// check VerifyToken if token valid or not
	VerifyToken(token string) (*Payload, error)
}
