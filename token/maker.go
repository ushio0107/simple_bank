// Since we need to implement both JWT and PASETO token makers, we can create a new interface called Maker.
// This interface will have two methods: CreateToken and VerifyToken. The JWTMaker struct will implement this interface.
package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
