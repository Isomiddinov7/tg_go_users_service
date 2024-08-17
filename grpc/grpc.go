package grpc

import (
	"tg_go_users_service/config"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/grpc/client"
	"tg_go_users_service/grpc/service"
	"tg_go_users_service/pkg/logger"
	"tg_go_users_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()
	users_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, srvc))
	users_service.RegisterUserSellOrBuyServiceServer(grpcServer, service.NewUserTransactionService(cfg, log, strg, srvc))
	users_service.RegisterUserMessageListServer(grpcServer, service.NewMessageService(cfg, log, strg, srvc))
	reflection.Register(grpcServer)
	return
}
