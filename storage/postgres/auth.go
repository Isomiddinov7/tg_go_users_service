package postgres

import (
	"context"
	"database/sql"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) storage.AuthRepoI {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) Auth(ctx context.Context, req *users_service.Req) (resp *users_service.AuthResp, err error) {
	var (
		query = `
			SELECT 
				id
			FROM "admin"
			WHERE "login" = $1 AND "password" = $2
		`

		id sql.NullString
	)
	err = r.db.QueryRow(ctx, query, req.Login, req.Password).Scan(
		&id,
	)
	if err != nil {
		return nil, err
	}

	return &users_service.AuthResp{
		Success: "Success",
		UserId:  id.String,
	}, nil
}

func (r *authRepo) Deserialize(ctx context.Context, req *users_service.DReq) (err error) {
	var (
		query = `
			SELECT 
				id
			FROM "admin"
			WHERE "id" = $1
		`

		id sql.NullString
	)
	err = r.db.QueryRow(ctx, query, req.AdminId).Scan(
		&id,
	)
	if err != nil {
		return err
	}

	return nil
}
