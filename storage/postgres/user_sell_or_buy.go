package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tg_go_users_service/genproto/coins_service"
	"tg_go_users_service/genproto/users_service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cast"
)

type userTransaction struct {
	db *pgxpool.Pool
}

func NewUserTransactionRepo(db *pgxpool.Pool) *userTransaction {
	return &userTransaction{
		db: db,
	}
}

func (r *userTransaction) UserSell(ctx context.Context, req *users_service.UserSellRequest) error {
	var (
		id    = uuid.New().String()
		query = `
			INSERT INTO "user_transaction"(
				id,
				user_id,
				coin_id,
				coin_amount,
				user_confirmation_img,
				coin_price,
				all_price,
				status,
				card_name,
				payment_card,
				message
			) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		coin_price sql.NullString
	)
	var (
		queryHalf = `
			SELECT
				"halfCoinAmount",
				"halfCoinPrice"
			FROM "half_coins_price"
			WHERE "coin_id" = $1
		`
	)

	rows, err := r.db.Query(ctx, queryHalf, req.CoinId)
	if err != nil {
		return err
	}
	halfPrices := []*coins_service.HalfCoinPrice{}
	for rows.Next() {
		var (
			halfPrice      = coins_service.HalfCoinPrice{}
			halfCoinAmount sql.NullString
			halfCoinPrice  sql.NullString
		)

		err = rows.Scan(
			&halfCoinAmount,
			&halfCoinPrice,
		)
		if err != nil {
			return err
		}
		halfPrice = coins_service.HalfCoinPrice{
			HalfCoinAmount: halfCoinAmount.String,
			HalfCoinPrice:  halfCoinPrice.String,
		}
		halfPrices = append(halfPrices, &halfPrice)
	}
	var (
		coin_sell_price sql.NullString
		summ            float64
	)
	if cast.ToInt64(req.CoinAmount) >= 1 {
		query1 := `
			SELECT
				"coin_sell_price"
			FROM "coins"
			WHERE id = $1
		`
		err = r.db.QueryRow(ctx, query1, req.CoinId).Scan(
			&coin_price,
		)
		if err != nil {
			return err
		}
		summ = cast.ToFloat64(req.CoinAmount) * cast.ToFloat64(coin_price.String)

	} else {
		for i := range halfPrices {
			if cast.ToFloat64(req.CoinAmount) == cast.ToFloat64(halfPrices[i].HalfCoinAmount) {
				coin_sell_price.String = cast.ToString(halfPrices[i].HalfCoinPrice)
			}
		}
		summ = cast.ToFloat64(coin_sell_price.String)
	}
	_, err = r.db.Exec(ctx, query,
		&id,
		req.UserId,
		req.CoinId,
		req.CoinAmount,
		req.CheckImg,
		&coin_price,
		cast.ToString(summ),
		"sell",
		req.CardHolderName,
		req.CardNumber,
		req.Message,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *userTransaction) UserBuy(ctx context.Context, req *users_service.UserBuyRequest) error {
	var (
		id    = uuid.New().String()
		query = `
			INSERT INTO "user_transaction"(
				id,
				coin_id,
				user_id,
				user_confirmation_img,
				coin_price,
				coin_amount,
				all_price,
				status,
				user_address,
				message
			) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		coin_price sql.NullString
	)

	var (
		queryHalf = `
			SELECT
				"halfCoinAmount",
				"halfCoinPrice"
			FROM "half_coins_price"
			WHERE "coin_id" = $1
		`
	)

	rows, err := r.db.Query(ctx, queryHalf, req.CoinId)
	if err != nil {
		return err
	}
	halfPrices := []*coins_service.HalfCoinPrice{}
	for rows.Next() {
		var (
			halfPrice      = coins_service.HalfCoinPrice{}
			halfCoinAmount sql.NullString
			halfCoinPrice  sql.NullString
		)

		err = rows.Scan(
			&halfCoinAmount,
			&halfCoinPrice,
		)
		if err != nil {
			return err
		}
		halfPrice = coins_service.HalfCoinPrice{
			HalfCoinAmount: halfCoinAmount.String,
			HalfCoinPrice:  halfCoinPrice.String,
		}
		halfPrices = append(halfPrices, &halfPrice)
	}
	var (
		coin_buy_price sql.NullString
		summ           float64
	)
	if cast.ToInt64(req.CoinAmount) >= 1 {
		query1 := `
			SELECT
				"coin_buy_price"
			FROM "coins"
			WHERE id = $1
		`
		err = r.db.QueryRow(ctx, query1, req.CoinId).Scan(
			&coin_price,
		)
		if err != nil {
			return err
		}
		summ = cast.ToFloat64(req.CoinAmount) * cast.ToFloat64(coin_price.String)

	} else {
		for i := range halfPrices {
			if cast.ToFloat64(req.CoinAmount) == cast.ToFloat64(halfPrices[i].HalfCoinAmount) {
				coin_buy_price.String = cast.ToString(halfPrices[i].HalfCoinPrice)
			}
		}
		summ = cast.ToFloat64(coin_buy_price.String)
	}
	_, err = r.db.Exec(ctx, query,
		&id,
		req.CoinId,
		req.UserId,
		req.PayImg,
		&coin_price,
		req.CoinAmount,
		cast.ToString(summ),
		"buy",
		req.Address,
		req.Message,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *userTransaction) AllUserSell(ctx context.Context, req *users_service.GetListUserTransactionRequest) (*users_service.GetListUserSellTransactionResponse, error) {

	var (
		query = `
			SELECT 
				COUNT(*) OVER(),
				ut.id,
				u.first_name,
				c.name,
				ut.coin_amount,
				ut.user_confirmation_img,
				ut.coin_price,
				ut.all_price,
				ut.status,
				ut.card_name,
				ut.payment_card,
				ut.message,
				ut.transaction_status,
				ut.created_at
			FROM "user_transaction" as ut
			JOIN "users" as u ON ut.user_id = u.id
			JOIN "coins" as c ON ut.coin_id = c.id
		`
		resp   users_service.GetListUserSellTransactionResponse
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
	if len(req.Status) > 0 {
		where = fmt.Sprintf(" WHERE ut.status = '%s'", req.Status)
	}

	query += where + sort + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id                    sql.NullString
			first_name            sql.NullString
			name                  sql.NullString
			coin_amount           sql.NullString
			user_confirmation_img sql.NullString
			coin_price            sql.NullString
			all_price             sql.NullString
			status                sql.NullString
			card_name             sql.NullString
			payment_card          sql.NullString
			message               sql.NullString
			transaction_status    sql.NullString
			created_at            sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&name,
			&coin_amount,
			&user_confirmation_img,
			&coin_price,
			&all_price,
			&status,
			&card_name,
			&payment_card,
			&message,
			&transaction_status,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		user_transaction := users_service.UserSellTransaction{
			Id:                id.String,
			UserId:            first_name.String,
			CoinId:            name.String,
			CoinAmount:        coin_amount.String,
			CheckImg:          user_confirmation_img.String,
			CoinPrice:         coin_price.String,
			AllPrice:          all_price.String,
			Status:            status.String,
			CardHolderName:    card_name.String,
			CardNumber:        payment_card.String,
			Message:           message.String,
			TransactionStatus: transaction_status.String,
			CreatedAt:         created_at.String,
		}

		resp.UserTransaction = append(resp.UserTransaction, &user_transaction)
	}

	return &resp, nil
}

func (r *userTransaction) AllUserBuy(ctx context.Context, req *users_service.GetListUserTransactionRequest) (*users_service.GetListUserBuyTransactionResponse, error) {
	var (
		query = `
			SELECT 
				COUNT(*) OVER(),
				ut.id,
				u.first_name,
				c.name,
				ut.coin_amount,
				ut.user_confirmation_img,
				ut.coin_price,
				ut.all_price,
				ut.status,
				ut.user_address,
				ut.message,
				ut.transaction_status,
				ut.created_at
			FROM "user_transaction" as ut
			JOIN "users" as u ON ut.user_id = u.id
			JOIN "coins" as c ON ut.coin_id = c.id
		`
		resp   users_service.GetListUserBuyTransactionResponse
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
	if len(req.Status) > 0 {
		where = fmt.Sprintf(" WHERE ut.status = '%s'", req.Status)
	}

	query += where + sort + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id                    sql.NullString
			first_name            sql.NullString
			name                  sql.NullString
			coin_amount           sql.NullString
			user_confirmation_img sql.NullString
			coin_price            sql.NullString
			all_price             sql.NullString
			status                sql.NullString
			user_address          sql.NullString
			message               sql.NullString
			transaction_status    sql.NullString
			created_at            sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&name,
			&coin_amount,
			&user_confirmation_img,
			&coin_price,
			&all_price,
			&status,
			&user_address,
			&message,
			&transaction_status,
			&created_at,
		)
		if err != nil {
			return nil, err
		}

		user_transaction := users_service.UserBuyTransaction{
			Id:                id.String,
			User:              first_name.String,
			Coin:              name.String,
			CoinAmount:        coin_amount.String,
			CheckImg:          user_confirmation_img.String,
			CoinPrice:         coin_price.String,
			AllPrice:          all_price.String,
			Status:            status.String,
			UserAddress:       user_address.String,
			Message:           message.String,
			TransactionStatus: transaction_status.String,
			CreatedAt:         created_at.String,
		}

		resp.UserTransaction = append(resp.UserTransaction, &user_transaction)
	}

	return &resp, nil
}

func (r *userTransaction) TransactionUpdate(ctx context.Context, req *users_service.UpdateTransaction) (int64, error) {
	var (
		query = `
			UPDATE "user_transaction"
			SET 
				"transaction_status" = $2,
				"updated_at" = NOW()
			WHERE "id" = $1
		`
	)

	rowsAffected, err := r.db.Exec(ctx,
		query,
		req.Id,
		req.TransactionStatus,
	)
	if err != nil {
		return 0, err
	}
	return rowsAffected.RowsAffected(), nil
}
