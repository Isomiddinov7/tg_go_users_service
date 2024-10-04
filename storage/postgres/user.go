package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cast"
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
				"telegram_id"
			) VALUES ($1, $2, $3, $4, $5)
		`
	)

	_, err := r.db.Exec(ctx, query,
		id,
		req.FirstName,
		req.LastName,
		req.Username,
		req.TelegramId,
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
				"status",
				"telegram_id",
				"created_at",
				"updated_at",
				"warnig_count",
				"block_time"
			FROM "users"
			WHERE "telegram_id"= $1
		`
		id           sql.NullString
		first_name   sql.NullString
		last_name    sql.NullString
		username     sql.NullString
		status       sql.NullString
		telegram_id  sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
		warnig_count sql.NullString
		block_time   sql.NullString
		user_status  users_service.UserStatus
	)

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&first_name,
		&last_name,
		&username,
		&status,
		&telegram_id,
		&created_at,
		&updated_at,
		&warnig_count,
		&block_time,
	)
	if err != nil {
		return nil, err
	}

	if cast.ToInt64(warnig_count.String) == 5 || cast.ToInt64(warnig_count.String) == 10 {
		user_status = users_service.UserStatus{
			Status:       block_time.String,
			WarningCount: warnig_count.String,
		}

	} else {
		user_status = users_service.UserStatus{
			Status:       status.String,
			WarningCount: warnig_count.String,
		}
	}

	return &users_service.User{
		Id:         id.String,
		FirstName:  first_name.String,
		LastName:   last_name.String,
		Username:   username.String,
		Status:     status.String,
		TelegramId: telegram_id.String,
		CreatedAt:  created_at.String,
		UpdatedAt:  updated_at.String,
		UserStatus: &user_status,
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
		where += " AND last_name ILIKE" + " '%" + req.Search + "%'" + " OR first_name ILIKE" + " '%" + req.Search + "%'" + " OR username ILIKE" + " '%" + req.Search + "%'"
	}

	query := `
		SELECT
			COUNT(*) OVER(),
			"id",
			"first_name",
			"last_name",
			"username",
			"status",
			"telegram_id",
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
			user        users_service.User
			id          sql.NullString
			first_name  sql.NullString
			last_name   sql.NullString
			username    sql.NullString
			status      sql.NullString
			telegram_id sql.NullString
			created_at  sql.NullString
			updated_at  sql.NullString
		)
		err = rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&last_name,
			&username,
			&status,
			&telegram_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		user = users_service.User{
			Id:         id.String,
			FirstName:  first_name.String,
			LastName:   last_name.String,
			Username:   username.String,
			Status:     status.String,
			TelegramId: telegram_id.String,
			CreatedAt:  created_at.String,
			UpdatedAt:  updated_at.String,
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
					"block_time" = $3,
					"updated_at" = NOW()
			WHERE "id" = $1
		`
		queryUpdate = `
			UPDATE "users"
				SET
					"status" = 'inactive',
					"block_time" = $2,
					"warnig_count" = $3,
					"updated_at" = NOW()
			WHERE "id" = $1
		`
		queryGet = `
			SELECT 
				"warnig_count"
			FROM "users"
			WHERE "id" = $1
		`
		warningCount sql.NullString
	)

	err := r.db.QueryRow(ctx, queryGet, req.Id).Scan(&warningCount)
	if err != nil {
		return 0, err
	}

	warning := cast.ToInt64(warningCount.String) + cast.ToInt64(req.WarningCount)

	if req.Status == "inactive" {
		blockTime := time.Now().Add(72 * time.Hour).Format(time.RFC3339)
		rowsAffected, err := r.db.Exec(ctx, query, req.Id, req.Status, blockTime)
		if err != nil {
			return 0, err
		}
		return rowsAffected.RowsAffected(), nil
	}

	var blockTime string
	if warning == 5 {
		blockTime = time.Now().Add(72 * time.Hour).Format(time.RFC3339)
	} else if warning == 10 {
		blockTime = time.Now().Add(360 * time.Hour).Format(time.RFC3339)
	}

	if blockTime != "" {
		rowsAffected, err := r.db.Exec(ctx, queryUpdate, req.Id, blockTime, warning)
		if err != nil {
			return 0, err
		}
		return rowsAffected.RowsAffected(), nil
	}

	return 0, nil
}
