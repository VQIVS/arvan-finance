package app

import (
	"billing-service/config"
	"billing-service/internal/user"
	userPort "billing-service/internal/user/port"
	"billing-service/pkg/adapters/rabbit"
	"billing-service/pkg/adapters/storage"
	"billing-service/pkg/postgres"
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type app struct {
	db          *gorm.DB
	cfg         config.Config
	userService userPort.Service
	rabbit      *rabbit.Rabbit
	logger      *slog.Logger
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Rabbit() *rabbit.Rabbit {
	return a.rabbit
}

func (a *app) UserService(ctx context.Context) userPort.Service {
	if a.userService == nil {
		a.userService = user.NewService(storage.NewUserRepo(a.db), a.rabbit)
	}
	return a.userService
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) Logger() *slog.Logger {
	return a.logger
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

	postgres.Migrate(db)

	a.db = db
	return nil
}
func NewApp(cfg config.Config) (App, error) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("trace_id", uuid.NewString())

	a := &app{
		cfg:    cfg,
		logger: l,
	}

	if err := a.setDB(); err != nil {
		return nil, err
	}
	if cfg.Rabbit.URL != "" {
		r, err := rabbit.NewRabbit(cfg.Rabbit.URL)
		if err != nil {
			return nil, err
		}
		a.rabbit = r
		if err := a.rabbit.InitQueues(cfg.Rabbit.Queues); err != nil {
			return nil, err
		}
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
