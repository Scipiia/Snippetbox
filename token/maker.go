package token

import (
	"time"
)

// интерфеис для создания, управления и проверки токена
type Maker interface {
	CreateToken(name string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
