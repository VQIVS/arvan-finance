package postgres

import (
	"billing-service/pkg/adapters/storage/types"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnOptions struct {
	User   string
	Pass   string
	Host   string
	Port   uint
	DBName string
	Schema string
}

func (o DBConnOptions) PostgresDSN() string {
	schema := o.Schema
	if schema == "" {
		schema = "public"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s",
		o.Host, o.Port, o.User, o.Pass, o.DBName, schema)
}

func NewPsqlGormConnection(opt DBConnOptions) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(opt.PostgresDSN()), &gorm.Config{
		Logger: logger.Discard,
	})
}

func Migrate(db *gorm.DB) {
	migrator := db.Migrator()
	// storage models
	if err := migrator.AutoMigrate(&types.User{}); err != nil {
		panic("failed to auto-migrate User table: " + err.Error())
	}
}
