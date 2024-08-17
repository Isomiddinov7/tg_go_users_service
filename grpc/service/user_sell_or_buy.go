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

type UserTransactionService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*users_service.UnimplementedUserSellOrBuyServiceServer
}

func NewUserTransactionService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *UserTransactionService {
	return &UserTransactionService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *UserTransactionService) UserSell(ctx context.Context, req *users_service.UserSellRequest) (*users_service.Empty, error) {

	i.log.Info("---UserSell------>", logger.Any("req", req))
	err := i.strg.UserTransaction().UserSell(ctx, req)
	if err != nil {
		i.log.Error("!!!UserTransaction->User->UserSell--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, nil
}

func (i *UserTransactionService) UserBuy(ctx context.Context, req *users_service.UserBuyRequest) (*users_service.Empty, error) {

	i.log.Info("---UserBuy------>", logger.Any("req", req))
	err := i.strg.UserTransaction().UserBuy(ctx, req)
	if err != nil {
		i.log.Error("!!!UserTransaction->User->UserBuy--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, nil
}

func (i *UserTransactionService) AllUserSell(ctx context.Context, req *users_service.GetListUserTransactionRequest) (resp *users_service.GetListUserSellTransactionResponse, err error) {

	i.log.Info("---AllUserSell------>", logger.Any("req", req))
	resp, err = i.strg.UserTransaction().AllUserSell(ctx, req)
	if err != nil {
		i.log.Error("!!!UserTransaction->User->AllUserSell--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *UserTransactionService) AllUserBuy(ctx context.Context, req *users_service.GetListUserTransactionRequest) (resp *users_service.GetListUserBuyTransactionResponse, err error) {

	i.log.Info("---AllUserBuy------>", logger.Any("req", req))
	resp, err = i.strg.UserTransaction().AllUserBuy(ctx, req)
	if err != nil {
		i.log.Error("!!!UserTransaction->User->AllUserBuy--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}
