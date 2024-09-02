package postgres

import (
	"context"
	"fmt"
	"tg_go_users_service/config"
	"tg_go_users_service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db               *pgxpool.Pool
	user             storage.UserRepoI
	user_transaction storage.UserTransactionRepoI
	message          storage.UserMessageRepoI
	auth             storage.AuthRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=require",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))

	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoI {
	if s.user == nil {
		s.user = NewUserRepo(s.db)
	}

	return s.user
}

func (s *Store) UserTransaction() storage.UserTransactionRepoI {
	if s.user_transaction == nil {
		s.user_transaction = NewUserTransactionRepo(s.db)
	}

	return s.user_transaction
}

func (s *Store) Messages() storage.UserMessageRepoI {
	if s.message == nil {
		s.message = NewUserMessageRepo(s.db)
	}

	return s.message
}

func (s *Store) Auth() storage.AuthRepoI {
	if s.auth == nil {
		s.auth = NewAuthRepo(s.db)
	}
	return s.auth
}
