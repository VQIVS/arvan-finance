package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	APIKey  string
	Balance float64
}
