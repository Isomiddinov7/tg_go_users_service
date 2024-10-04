package service

import (
	"context"
	"tg_go_users_service/config"
	"tg_go_users_service/genproto/users_service"
	"tg_go_users_service/grpc/client"
	"tg_go_users_service/pkg/logger"
	"tg_go_users_service/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	users_service.UnimplementedAuthServiceServer
}

func NewAuthService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *AuthService {
	return &AuthService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *AuthService) Auth(ctx context.Context, req *users_service.Req) (resp *users_service.AuthResp, err error) {
	resp, err = i.strg.Auth().Auth(ctx, req)
	if err != nil {
		i.log.Error("!!!user service->Login->Login--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}
func (i *AuthService) Deserialize(ctx context.Context, req *users_service.DReq) (resp *users_service.Empty, err error) {
	err = i.strg.Auth().Deserialize(ctx, req)
	if err != nil {
		i.log.Error("!!!user service->Deserialize->Deserialize--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}
