package app

import (
	"context"
	"finance/config"
	"finance/internal/infra/messaging"
	"finance/internal/infra/storage"
	"finance/internal/infra/storage/types"
	"finance/internal/usecase"
	"finance/pkg/logger"
	"finance/pkg/postgres"
	"finance/pkg/rabbit"

	"gorm.io/gorm"
)

type app struct {
	db            *gorm.DB
	cfg           config.Config
	rabbitConn    *rabbit.RabbitConn
	walletService *usecase.WalletService
	logger        *logger.Logger
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) RabbitConn() *rabbit.RabbitConn {
	return a.rabbitConn
}

func (a *app) WalletService(ctx context.Context) *usecase.WalletService {
	return a.walletService
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg:    cfg,
		logger: logger.NewLogger(""),
	}
	if err := a.setDB(); err != nil {
		return nil, err
	}

	if err := a.setRabbitConn(); err != nil {
		return nil, err
	}

	if err := a.initQueues(); err != nil {
		return nil, err
	}

	a.setWalletService(a.db, a.rabbitConn)
	return a, nil
}

func NewMustApp(cfg config.Config) App {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return app
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
	// Auto migrate
	err = postgres.Migrate(db, &types.Wallet{}, &types.Transaction{}, &types.User{})
	if err != nil {
		return err
	}

	a.db = db
	return nil
}

func (a *app) setRabbitConn() error {
	rabbitConn := rabbit.NewRabbitConn(a.cfg.RabbitMQ.URI)
	a.rabbitConn = rabbitConn
	return nil
}

func (a *app) initQueues() error {
	for _, q := range a.cfg.RabbitMQ.Queues {
		err := a.rabbitConn.DeclareBindQueue(q.Name, q.Exchange, q.Routing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) setWalletService(db *gorm.DB, rabbitConn *rabbit.RabbitConn) {
	walletRepo := storage.NewWalletRepository(db)
	transactionRepo := storage.NewTransactionRepo(db)
	userRepo := storage.NewUserRepository(db)
	txManager := storage.NewGormTransactionManager(db)
	walletPublisher := messaging.NewWalletPublisher(rabbitConn, a.logger)
	a.walletService = usecase.NewWalletService(walletRepo, userRepo, transactionRepo, txManager, walletPublisher, a.logger)
}
