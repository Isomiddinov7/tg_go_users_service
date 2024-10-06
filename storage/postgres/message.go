package postgres

import (
	"context"
	"database/sql"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cast"
)

type userMessageRepo struct {
	db *pgxpool.Pool
}

func NewUserMessageRepo(db *pgxpool.Pool) storage.UserMessageRepoI {
	return &userMessageRepo{
		db: db,
	}
}

func (r *userMessageRepo) CreateUserMessage(ctx context.Context, req *users_service.MessageRequest) error {
	var (
		messageId = uuid.New().String()
		query     = `
			INSERT INTO "messages"(
				id,
				status,
				message,
				read,
				admin_id,
				user_id,
				file
			) VALUES($1, $2, $3, $4, $5, $6, $7)
		`
	)

	_, err := r.db.Exec(ctx, query,
		&messageId,
		"user",
		req.Text,
		"false",
		"dbecf401-64b3-4b9b-829a-c8b061431286",
		req.UserId,
		req.File,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userMessageRepo) CreateAdminMessage(ctx context.Context, req *users_service.MessageRequest) error {

	var (
		messageId = uuid.New().String()

		query = `
			INSERT INTO "messages"(
				id,
				status,
				message,
				read,
				admin_id,
				user_id,
				file
			) VALUES($1, $2, $3, $4, $5, $6, $7)
		`
	)

	_, err := r.db.Exec(ctx, query,
		messageId,
		"admin",
		req.Text,
		"false",
		"dbecf401-64b3-4b9b-829a-c8b061431286",
		req.UserId,
		req.File,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *userMessageRepo) UpdateMessage(ctx context.Context, req *users_service.UpdateMessageRequest) (int64, error) {
	var (
		query = `
			UPDATE "messages"
			SET 
				"read" = $2,
				"updated_at" = NOW()
			WHERE "id" = $1
		`
	)

	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.Read,
	)
	if err != nil {
		return 0, err
	}
	return rowsAffected.RowsAffected(), nil
}

func (r *userMessageRepo) GetUserMessage(ctx context.Context, req *users_service.GetMessageUserRequest) (*users_service.GetMessageUserResponse, error) {
	var (
		query = `
			SELECT
				COUNT(*) OVER(),
				id,
				status,
				message,
				read,
				user_id,
				file,
				created_at
			FROM "messages"
			WHERE "user_id" = $1
		`
		resp users_service.GetMessageUserResponse
	)

	rows, err := r.db.Query(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			message    users_service.Message
			id         sql.NullString
			status     sql.NullString
			text       sql.NullString
			read       sql.NullString
			user_id    sql.NullString
			file       sql.NullString
			created_at sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&status,
			&text,
			&read,
			&user_id,
			&file,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		message = users_service.Message{
			Id:        id.String,
			Text:      text.String,
			File:      file.String,
			Status:    status.String,
			Read:      read.String,
			UserId:    user_id.String,
			CreatedAt: created_at.String,
		}
		resp.Messages = append(resp.Messages, &message)
	}

	return &resp, nil
}

func (r *userMessageRepo) GetAdminAllMessage(ctx context.Context) (*users_service.GetMessageAdminResponse, error) {
	var (
		query = `
			WITH user_messages AS (
				SELECT DISTINCT ON (m.user_id)
					m.user_id,
					u.first_name,
					u.last_name,
					m.file,
					m.message,
					m.read,
					m.status,
					m.created_at,
					m.updated_at
				FROM "messages" as m
				JOIN "users" as u ON u.id = m.user_id
				ORDER BY m.user_id, m.created_at DESC
			)
			SELECT
				(SELECT COUNT(DISTINCT user_id) FROM "messages" WHERE "status" = 'user') as user_count,
				user_messages.user_id,
				user_messages.first_name,
				user_messages.last_name,
				user_messages.file,
				user_messages.message,
				user_messages.read,
				user_messages.status,
				user_messages.created_at,
				user_messages.updated_at
			FROM user_messages
		`
		resp            users_service.GetMessageAdminResponse
		userMessagesMap = make(map[string]*users_service.AdminResponse)
	)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			message_response_message users_service.AdminResponseMessage
			first_name               sql.NullString
			last_name                sql.NullString
			user_id                  sql.NullString
			file                     sql.NullString
			message                  sql.NullString
			read                     sql.NullString
			status                   sql.NullString
			created_at               sql.NullString
			updated_at               sql.NullString
		)
		err = rows.Scan(
			&resp.MessageCount,
			&user_id,
			&first_name,
			&last_name,
			&file,
			&message,
			&read,
			&status,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		message_response_message = users_service.AdminResponseMessage{
			LastMessage: message.String,
			File:        file.String,
			Read:        read.String,
			Status:      status.String,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		}

		userMessagesMap[user_id.String] = &users_service.AdminResponse{
			UserId:    user_id.String,
			FirstName: first_name.String,
			LastName:  last_name.String,
			Message:   []*users_service.AdminResponseMessage{&message_response_message},
		}
	}

	for _, messageResponse := range userMessagesMap {
		resp.AdminMessage = append(resp.AdminMessage, messageResponse)
	}

	return &resp, nil
}

// func (r *userMessageRepo) GetAdminAllMessage(ctx context.Context) (*users_service.GetMessageAdminResponse, error) {
// 	var (
// 		query = `
// 			SELECT
// 				COUNT(m.*) OVER(),
// 				m.user_id,
// 				u.first_name,
// 				u.last_name,
// 				m.file,
// 				m.message,
// 				m.read,
// 				m.created_at,
// 				m.updated_at
// 			FROM "messages" as m
// 			JOIN "users" as u ON u.id=m.user_id
// 			WHERE m."status" = 'user'
// 			ORDER BY m.created_at ASC
// 		`
// 		resp            users_service.GetMessageAdminResponse
// 		userMessagesMap = make(map[string]*users_service.AdminResponse)
// 	)
// 	rows, err := r.db.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var (
// 			message_response_message users_service.AdminResponseMessage
// 			last_message             sql.NullString
// 			first_name               sql.NullString
// 			last_name                sql.NullString
// 			user_id                  sql.NullString
// 			file                     sql.NullString
// 			read                     sql.NullString
// 			created_at               sql.NullString
// 			updated_at               sql.NullString
// 		)
// 		err = rows.Scan(
// 			&resp.MessageCount,
// 			&user_id,
// 			&first_name,
// 			&last_name,
// 			&file,
// 			&last_message,
// 			&read,
// 			&created_at,
// 			&updated_at,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		message_response_message = users_service.AdminResponseMessage{
// 			LastMessage: last_message.String,
// 			File:        file.String,
// 			Read:        read.String,
// 			CreatedAt:   created_at.String,
// 			UpdatedAt:   updated_at.String,
// 		}

// 		if userResponse, exists := userMessagesMap[user_id.String]; exists {
// 			userResponse.Message = append(userResponse.Message, &message_response_message)
// 		} else {
// 			userMessagesMap[user_id.String] = &users_service.AdminResponse{
// 				UserId:    user_id.String,
// 				FirstName: first_name.String,
// 				LastName:  last_name.String,
// 				Message:   []*users_service.AdminResponseMessage{&message_response_message},
// 			}
// 		}
// 	}

// 	for _, messageResponse := range userMessagesMap {
// 		resp.AdminMessage = append(resp.AdminMessage, messageResponse)
// 	}

// 	return &resp, nil
// }

func (r *userMessageRepo) GetMessageAdminID(ctx context.Context, req *users_service.GetMessageUserRequest) (*users_service.GetMessageAdminById, error) {
	var (
		query = `
			SELECT
				u.first_name,
				u.last_name,
				m.file,
				m.id,
				m.message,
				m.status,
				m.read,
				m.user_id,
				m.created_at
			FROM "messages" as m
			JOIN "users" as u ON u.id = m.user_id
			WHERE m.user_id = $1
			ORDER BY m.created_at DESC

		`
		first_name sql.NullString
		last_name  sql.NullString
		resp       users_service.GetMessageAdminById
	)
	rows, err := r.db.Query(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			messages   users_service.Message
			id         sql.NullString
			message    sql.NullString
			file       sql.NullString
			status     sql.NullString
			read       sql.NullString
			user_id    sql.NullString
			created_at sql.NullString
		)
		err = rows.Scan(
			&first_name,
			&last_name,
			&file,
			&id,
			&message,
			&status,
			&read,
			&user_id,
			&created_at,
		)
		if err != nil {
			return nil, err
		}
		messages = users_service.Message{
			Id:        id.String,
			Text:      message.String,
			File:      file.String,
			Status:    status.String,
			Read:      read.String,
			UserId:    user_id.String,
			CreatedAt: created_at.String,
		}

		resp.Message = append(resp.Message, &messages)
		resp.FirstName = first_name.String
		resp.LastName = last_name.String
		resp.File = file.String
	}
	return &resp, nil
}

func (r *userMessageRepo) SendMessageUser(ctx context.Context, req *users_service.TelegramMessageUser) (resp *users_service.TelegramMessageResponse, err error) {
	var (
		query = `
			SELECT
				m."message",
				u."telegram_id",
				m."file"
				m.created_at,
				m.updated_at
			FROM "messages" as m
			JOIN "users" as u ON u."id"="user_id"
			WHERE m.status = 'admin' AND m.user_id = $1 
			ORDER BY m.created_at DESC
			LIMIT 1
		`
		message     sql.NullString
		telegram_id sql.NullString
		file        sql.NullString
		created_at  sql.NullString
		updated_at  sql.NullString
	)

	err = r.db.QueryRow(ctx, query, req.Id).Scan(
		&message,
		&telegram_id,
		&file,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, err
	}

	resp = &users_service.TelegramMessageResponse{
		Message:    message.String,
		TelegramId: telegram_id.String,
		File:       file.String,
		CreatedAt:  created_at.String,
		UpdatedAt:  updated_at.String,
	}
	return resp, nil

}

func (r *userMessageRepo) PayMessagePost(ctx context.Context, req *users_service.PaymsqRequest) error {
	var (
		query = `
			INSERT INTO "pay_message"(
				"id",
				"message",
				"file",
				"user_id",
				"user_transaction_id",
				"premium_transaction_id"
			) VALUES($1, $2, $3, $4, $5, $6)

		`
		id = uuid.NewString()

		queryUser = `
			UPDATE "user_transaction" 
				SET 
					"transaction_status" = 'success',
					"updated_at" = NOW()
			WHERE id = $1
		`

		queryPremium = `
			UPDATE "premium_transaction"
				SET 
					"transaction_status" = 'success',
					"updated_at" = NOW()
			WHERE id = $1
		`

		queryUserError = `
			UPDATE "user_transaction" 
				SET 
					"transaction_status" = 'error',
					"updated_at" = NOW()
			WHERE id = $1
		`

		queryPremiumError = `
			UPDATE "premium_transaction"
				SET 
					"transaction_status" = 'error',
					"updated_at" = NOW()
			WHERE id = $1
		`

		queryUserPut = `
			UPDATE "users"
				SET 
					"warnig_count" = $2,
					"status" = $3,
					"block_time" = $4
			WHERE id = $1
		`
		queryUserGet = `
			SELECT 
				"warnig_count"
			FROM "users"
			WHERE "id" = $1
		`
		warnig_count sql.NullString
	)

	_, err := r.db.Exec(ctx, query,
		id,
		req.Message,
		req.File,
		req.UserId,
		req.UserTransactionId,
		req.PremiumTransactionId,
	)
	if err != nil {
		return err
	}
	if len(req.PremiumTransactionId) > 0 && len(req.Status) == 0 {
		_, err = r.db.Exec(ctx, queryPremium, req.PremiumTransactionId)
		if err != nil {
			return err
		}
	} else if len(req.PremiumTransactionId) > 0 && len(req.Status) > 0 {
		_, err := r.db.Exec(ctx, queryPremiumError, req.PremiumTransactionId)
		if err != nil {
			return err
		}
	}

	if len(req.UserTransactionId) > 0 && len(req.Status) == 0 {
		_, err = r.db.Exec(ctx, queryUser, req.UserTransactionId)
		if err != nil {
			return err
		}
	} else if len(req.UserTransactionId) > 0 && len(req.Status) > 0 {
		var (
			sum    int64
			status = "active"
		)
		_, err = r.db.Exec(ctx, queryUserError, req.UserTransactionId)
		if err != nil {
			return err
		}

		err = r.db.QueryRow(ctx, queryUserGet, req.UserId).Scan(
			&warnig_count,
		)
		if err != nil {
			return err
		}

		sum = cast.ToInt64(warnig_count.String) + 1
		var (
			block_time time.Time
		)
		if sum%5 == 0 {
			block_time = time.Now().Add(3 * time.Hour)
			status = "inactive"
		}

		_, err = r.db.Exec(ctx, queryUserPut, req.UserId, cast.ToString(sum), &status, block_time)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *userMessageRepo) PayMessageGet(ctx context.Context, req *users_service.PaymsqUser) (*users_service.PaymsqResponse, error) {
	var (
		query = `
			SELECT 
				"id",
				"file",
				"message",
				"user_transaction_id",
				"premium_transaction_id",
				"created_at",
				"updated_at"
			FROM "pay_message"
			WHERE "user_id" = $1

		`
		resp                   = &users_service.PaymsqResponse{}
		messages               []*users_service.Paymsq
		id                     sql.NullString
		file                   sql.NullString
		message                sql.NullString
		user_transaction_id    sql.NullString
		premium_transaction_id sql.NullString
		created_at             sql.NullString
		updated_at             sql.NullString
	)

	rows, err := r.db.Query(ctx, query, req.UserTransactionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(
			&id,
			&file,
			&message,
			&user_transaction_id,
			&premium_transaction_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &users_service.Paymsq{
			Id:                   id.String,
			File:                 file.String,
			Message:              message.String,
			UserTransactionId:    user_transaction_id.String,
			PremiumTransactionId: premium_transaction_id.String,
			CreatedAt:            created_at.String,
			UpdatedAt:            updated_at.String,
		})

		resp.Message = messages
	}

	return resp, nil
}
