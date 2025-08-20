package types

import (
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	UserID uint
	User   *User
	Amount float64
	SMSID  uint
}
