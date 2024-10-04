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

type MessageService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	users_service.UnimplementedUserMessageListServer
}

func NewMessageService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *MessageService {
	return &MessageService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (i *MessageService) CreateUserMessage(ctx context.Context, req *users_service.MessageRequest) (resp *users_service.Empty, err error) {

	i.log.Info("---CreateUserMessage------>", logger.Any("req", req))
	err = i.strg.Messages().CreateUserMessage(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateUserMessage->User->CreateMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MessageService) CreateAdminMessage(ctx context.Context, req *users_service.MessageRequest) (resp *users_service.Empty, err error) {

	i.log.Info("---CreateAdminMessage------>", logger.Any("req", req))
	err = i.strg.Messages().CreateAdminMessage(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateAdminMessage->Admin->CreateMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MessageService) UpdateMessage(ctx context.Context, req *users_service.UpdateMessageRequest) (resp *users_service.Empty, err error) {

	i.log.Info("---UpdateMessage------>", logger.Any("req", req))
	_, err = i.strg.Messages().UpdateMessage(ctx, req)
	if err != nil {
		i.log.Error("!!!UpdateMessage->Message->UpdateMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (i *MessageService) GetUserMessage(ctx context.Context, req *users_service.GetMessageUserRequest) (resp *users_service.GetMessageUserResponse, err error) {

	i.log.Info("---GetUserMessage------>", logger.Any("req", ""))
	resp, err = i.strg.Messages().GetUserMessage(ctx, req)
	if err != nil {
		i.log.Error("!!!GetUserMessage->Message->GetUserMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (i *MessageService) GetAdminAllMessage(ctx context.Context, req *users_service.Empty) (resp *users_service.GetMessageAdminResponse, err error) {
	i.log.Info("---GetAdminAllMessage------>", logger.Any("req", ""))
	resp, err = i.strg.Messages().GetAdminAllMessage(ctx)
	if err != nil {
		i.log.Error("!!!GetAdminAllMessage->Message->GetAdminAllMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func (i *MessageService) GetMessageAdminID(ctx context.Context, req *users_service.GetMessageUserRequest) (resp *users_service.GetMessageAdminById, err error) {
	i.log.Info("---GetAdminAllMessage------>", logger.Any("req", ""))
	resp, err = i.strg.Messages().GetMessageAdminID(ctx, req)
	if err != nil {
		i.log.Error("!!!GetAdminAllMessage->Message->GetAdminAllMessage--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func (i *MessageService) SendMessageUser(ctx context.Context, req *users_service.TelegramMessageUser) (resp *users_service.TelegramMessageResponse, err error) {
	i.log.Info("---SendMessageUser------>", logger.Any("req", ""))
	resp, err = i.strg.Messages().SendMessageUser(ctx, req)
	if err != nil {
		i.log.Error("!!!SendMessageUser->Message->SendMessageUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func (i *MessageService) PayMessagePost(ctx context.Context, req *users_service.PaymsqRequest) (resp *users_service.Empty, err error) {
	i.log.Info("---PayMessagePost------>", logger.Any("req", req))
	err = i.strg.Messages().PayMessagePost(ctx, req)
	if err != nil {
		i.log.Error("!!!PayMessagePost->Message->PayMessagePost--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func (i *MessageService) PayMessageGet(ctx context.Context, req *users_service.PaymsqUser) (resp *users_service.PaymsqResponse, err error) {
	i.log.Info("---PayMessageGet------>", logger.Any("req", ""))
	resp, err = i.strg.Messages().PayMessageGet(ctx, req)
	if err != nil {
		i.log.Error("!!!PayMessageGet->Message->PayMessageGet--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}
