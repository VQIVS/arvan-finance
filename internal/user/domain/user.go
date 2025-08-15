package domain

import "time"

type (
	UserID uint
	APIKey string
)

type User struct {
	ID        UserID
	APIKey    APIKey
	Balance   float64
	CreatedAt time.Time
	DeletedAt time.Time
}

func (u *User) HasSufficientBalance(amount uint64) bool {
	return u.Balance >= float64(amount)
}
