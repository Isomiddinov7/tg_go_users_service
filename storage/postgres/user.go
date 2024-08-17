package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserRepoI {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, req *users_service.CreateUser) error {
	var (
		id    = uuid.NewString()
		query = `
			INSERT INTO "users"(
				"id",
				"first_name",
				"last_name",
				"username",
				"auth_date",
				"hash",
				"status"
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
	)

	_, err := r.db.Exec(ctx, query,
		id,
		req.FirstName,
		req.LastName,
		req.Username,
		req.AuthDate,
		req.Hash,
		req.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, req *users_service.UserPrimaryKey) (*users_service.User, error) {
	var (
		query = `
			SELECT
				"id",
				"first_name",
				"last_name",
				"username",
				"auth_date",
				"hash",
				"status",
				"created_at",
				"updated_at"
			FROM "users"
			WHERE "id"= $1
		`
		id         sql.NullString
		first_name sql.NullString
		last_name  sql.NullString
		username   sql.NullString
		auth_date  sql.NullString
		hash       sql.NullString
		status     sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&first_name,
		&last_name,
		&username,
		&auth_date,
		&hash,
		&status,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, err
	}

	return &users_service.User{
		Id:        id.String,
		FirstName: first_name.String,
		LastName:  last_name.String,
		Username:  username.String,
		AuthDate:  auth_date.String,
		Hash:      hash.String,
		Status:    status.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (r *userRepo) GetAll(ctx context.Context, req *users_service.GetListUserRequest) (*users_service.GetListUserResponse, error) {
	var (
		resp   users_service.GetListUserResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if len(req.Search) > 0 {
		where += " AND user_name ILIKE" + " '%" + req.Search + "%'"
	}

	query := `
		SELECT
			COUNT(*) OVER(),
			"id",
			"first_name",
			"last_name",
			"username",
			"auth_date",
			"hash",
			"status",
			"created_at",
			"updated_at"
		FROM "users"
	`

	query += where + sort + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			user       users_service.User
			id         sql.NullString
			first_name sql.NullString
			last_name  sql.NullString
			username   sql.NullString
			auth_date  sql.NullString
			hash       sql.NullString
			status     sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)
		err = rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&last_name,
			&username,
			&auth_date,
			&hash,
			&status,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		user = users_service.User{
			Id:        id.String,
			FirstName: first_name.String,
			LastName:  last_name.String,
			Username:  username.String,
			AuthDate:  auth_date.String,
			Hash:      hash.String,
			Status:    status.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		}

		resp.Users = append(resp.Users, &user)
	}

	return &resp, nil
}

func (r *userRepo) Update(ctx context.Context, req *users_service.UpdateUser) (int64, error) {
	var (
		query = `
			UPDATE "users"
				SET
					"status" = $2,
					"updated_at" = NOW()
			WHERE "id" = $1
		`
	)

	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Status,
	)

	if err != nil {
		return 0, err
	}
	return rowsAffected.RowsAffected(), nil
}
