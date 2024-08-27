package storage

import (
	"context"
	"tg_go_users_service/genproto/users_service"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
	UserTransaction() UserTransactionRepoI
	Messages() UserMessageRepoI
	Auth() AuthRepoI
}

type UserRepoI interface {
	Create(ctx context.Context, req *users_service.CreateUser) error
	GetByID(ctx context.Context, req *users_service.UserPrimaryKey) (resp *users_service.User, err error)
	GetAll(ctx context.Context, req *users_service.GetListUserRequest) (resp *users_service.GetListUserResponse, err error)
	Update(ctx context.Context, req *users_service.UpdateUser) (rowsAffected int64, err error)
}

type UserMessageRepoI interface {
	CreateUserMessage(ctx context.Context, req *users_service.MessageRequest) error
	CreateAdminMessage(ctx context.Context, req *users_service.MessageRequest) error
	UpdateMessage(ctx context.Context, req *users_service.UpdateMessageRequest) (int64, error)
	GetUserMessage(ctx context.Context, req *users_service.GetMessageUserRequest) (resp *users_service.GetMessageUserResponse, err error)
	GetAdminAllMessage(ctx context.Context) (resp *users_service.GetMessageAdminResponse, err error)
	GetMessageAdminID(ctx context.Context, req *users_service.GetMessageUserRequest) (resp *users_service.GetMessageAdminById, err error)
}

type UserTransactionRepoI interface {
	UserSell(ctx context.Context, req *users_service.UserSellRequest) error
	UserBuy(ctx context.Context, req *users_service.UserBuyRequest) error
	AllUserSell(ctx context.Context, req *users_service.GetListUserTransactionRequest) (*users_service.GetListUserSellTransactionResponse, error)
	AllUserBuy(ctx context.Context, req *users_service.GetListUserTransactionRequest) (*users_service.GetListUserBuyTransactionResponse, error)
}

type AuthRepoI interface {
	SignIn(ctx context.Context, req *users_service.Credentials) (resp *users_service.CredetialId, err error)
}
