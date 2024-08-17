package client

import (
	"tg_go_users_service/config"
	"tg_go_users_service/genproto/users_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManagerI interface {
	UserService() users_service.UserServiceClient
	UserTransactionService() users_service.UserSellOrBuyServiceClient
	Messages() users_service.UserMessageListClient
}

type grpcClients struct {
	userService            users_service.UserServiceClient
	usertransactionService users_service.UserSellOrBuyServiceClient
	messageService         users_service.UserMessageListClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {
	connUsersService, err := grpc.NewClient(
		cfg.UsersServiceHost+cfg.UsersGRPCPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &grpcClients{
		userService:            users_service.NewUserServiceClient(connUsersService),
		usertransactionService: users_service.NewUserSellOrBuyServiceClient(connUsersService),
		messageService:         users_service.NewUserMessageListClient(connUsersService),
	}, nil
}

func (g *grpcClients) UserService() users_service.UserServiceClient {
	return g.userService
}

func (g *grpcClients) UserTransactionService() users_service.UserSellOrBuyServiceClient {
	return g.usertransactionService
}

func (g *grpcClients) Messages() users_service.UserMessageListClient {
	return g.messageService
}
