package main

import (
	"finance/config"
	"finance/internal/infra/storage/types"
	"finance/pkg/postgres"
	"fmt"
	"log"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   cfg.DB.User,
		Pass:   cfg.DB.Password,
		Host:   cfg.DB.Host,
		Port:   cfg.DB.Port,
		DBName: cfg.DB.Database,
		Schema: cfg.DB.Schema,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Starting database seeding...")

	if err := clearData(db); err != nil {
		log.Fatalf("Failed to clear existing data: %v", err)
	}

	if err := seedUsersAndWallets(db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	fmt.Println("Database seeding completed successfully!")
}

func clearData(db *gorm.DB) error {
	fmt.Println("Clearing existing data...")

	if err := db.Exec("DELETE FROM transactions").Error; err != nil {
		return fmt.Errorf("failed to clear transactions: %w", err)
	}

	if err := db.Exec("DELETE FROM wallets").Error; err != nil {
		return fmt.Errorf("failed to clear wallets: %w", err)
	}

	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	fmt.Println("Existing data cleared successfully")
	return nil
}

func seedUsersAndWallets(db *gorm.DB) error {
	fmt.Println("Seeding users and wallets...")

	usersData := []struct {
		Name     string
		LastName string
		Phone    string
		Balance  int64
		Currency string
	}{
		{"Ahmad", "Ahmadi", "+989123456001", 100000, "IRR"},
		{"Sara", "Saravi", "+989123456002", 250000, "IRR"},
		{"Mohammad", "Mohammadi", "+989123456003", 50000, "IRR"},
		{"Fateme", "Fatemi", "+989123456004", 300000, "IRR"},
		{"Ali", "Alavi", "+989123456005", 150000, "IRR"},
		{"Zahra", "Zahravi", "+989123456006", 200000, "IRR"},
		{"Hassan", "Hassani", "+989123456007", 75000, "IRR"},
		{"Maryam", "Maryami", "+989123456008", 400000, "IRR"},
		{"Reza", "Rezaei", "+989123456009", 120000, "IRR"},
		{"Neda", "Nedavi", "+989123456010", 180000, "IRR"},
	}

	for i, userData := range usersData {
		// Create user
		user := types.User{
			Base: types.Base{
				ID: uuid.New(),
			},
			Name:     userData.Name,
			LastName: userData.LastName,
			Phone:    userData.Phone,
			WalletID: nil,
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", userData.Name, err)
		}

		wallet := types.Wallet{
			Base: types.Base{
				ID: uuid.New(),
			},
			UserID:   user.ID,
			Balance:  types.NewBigInt(big.NewInt(userData.Balance)),
			Currency: userData.Currency,
		}

		if err := db.Create(&wallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet for user %s: %w", userData.Name, err)
		}

		walletIDStr := wallet.ID.String()
		user.WalletID = &walletIDStr
		if err := db.Save(&user).Error; err != nil {
			return fmt.Errorf("failed to update user %s with wallet ID: %w", userData.Name, err)
		}

		fmt.Printf("Created user %d: %s %s (Phone: %s) with wallet balance: %d %s\n",
			i+1, userData.Name, userData.LastName, userData.Phone, userData.Balance, userData.Currency)
	}

	return nil
}
