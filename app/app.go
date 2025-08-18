package app

import (
	"billing-service/config"
	"billing-service/internal/user"
	userPort "billing-service/internal/user/port"
	"context"

	"billing-service/pkg/adapters/rabbit"
	"billing-service/pkg/adapters/storage"
	appCtx "billing-service/pkg/context"
	"billing-service/pkg/postgres"

	"gorm.io/gorm"
)

type app struct {
	db          *gorm.DB
	cfg         config.Config
	userService userPort.Service
	rabbit      *rabbit.Rabbit
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Rabbit() *rabbit.Rabbit {
	return a.rabbit
}

func (a *app) UserService(ctx context.Context) userPort.Service {
	db := appCtx.GetDB(ctx)
	if db == nil {
		if a.userService == nil {
			a.userService = a.userServiceWithDB(a.db)
		}
		return a.userService
	}

	return a.userServiceWithDB(db)
}

func (a *app) userServiceWithDB(db *gorm.DB) userPort.Service {
	return user.NewService(storage.NewUserRepo(db), a.rabbit)
}

func (a *app) Config() config.Config {
	return a.cfg
}
func (a *app) setDB() error {
	db, err := postgres.NewPsqlGormConnection(postgres.DBConnOptions{
		User:   a.cfg.DB.User,
		Pass:   a.cfg.DB.Password,
		Host:   a.cfg.DB.Host,
		Port:   a.cfg.DB.Port,
		DBName: a.cfg.DB.Database,
		Schema: a.cfg.DB.Schema,
	})

	if err != nil {
		return err
	}

	// run auto migrations for storage models
	postgres.Migrate(db)

	a.db = db
	return nil
}
func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}
	// initialize rabbit connection if configured
	if cfg.Rabbit.URL != "" {
		r, err := rabbit.NewRabbit(cfg.Rabbit.URL)
		if err != nil {
			return nil, err
		}
		a.rabbit = r
	}
	return a, nil

}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
}
