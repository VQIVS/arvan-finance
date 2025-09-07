package types

type User struct {
	Base
	Name     string  `gorm:"type:varchar(100);index"`
	LastName string  `gorm:"type:varchar(100);index"`
	Phone    string  `gorm:"type:varchar(20);uniqueIndex"`
	WalletID *string `gorm:"type:varchar(36);index"`
}
