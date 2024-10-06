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

type UserService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	users_service.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *UserService {
	return &UserService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *UserService) GetById(ctx context.Context, req *users_service.UserPrimaryKey) (resp *users_service.User, err error) {
	i.log.Info("---GetUserByID------>", logger.Any("req", req))

	resp, err = i.strg.User().GetByID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetUserByID->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func (i *UserService) GetList(ctx context.Context, req *users_service.GetListUserRequest) (resp *users_service.GetListUserResponse, err error) {

	i.log.Info("---GetUsers------>", logger.Any("req", req))
	resp, err = i.strg.User().GetAll(ctx, req)
	if err != nil {
		i.log.Error("!!!GetUsers->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (i *UserService) Create(ctx context.Context, req *users_service.CreateUser) (resp *users_service.User, err error) {
	i.log.Info("---CreateUser------>", logger.Any("req", req))
	err = i.strg.User().Create(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateUser->User->TelegramStart--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (i *UserService) Update(ctx context.Context, req *users_service.UpdateUserStatus) (resp *users_service.Empty, err error) {
	i.log.Info("---UpdateUser------>", logger.Any("req", req))
	err = i.strg.User().Update(ctx, req)
	if err != nil {
		i.log.Error("!!!UpdateUser->User->Telegram_active--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}
