package types

type User struct {
	Base
	Name     string
	LastName string
	Phone    string
	WalletID *string
	Wallet   *Wallet
}
