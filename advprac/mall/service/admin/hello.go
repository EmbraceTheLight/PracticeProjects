package admin

import (
	"context"
	"go.uber.org/zap"
	"mall/common"
	"mall/service/do"
	"mall/service/dto"
	"mall/utils/logger"
)

func (s *Service) HelloWorld(ctx context.Context, adminUser *common.AdminUser, req *dto.HelloWorldReq) (*dto.HelloWorldResp, common.Errno) {
	resp, err := s.adminUSer.HelloWorld(ctx, &do.HelloWorld{})
	if err != nil {

		logger.Error("HelloWorld error", zap.Error(err), zap.Any("req", req))
		return nil, common.DatabaseErr.WithError(err)
	}
	return &dto.HelloWorldResp{
		Hello: "hello: " + resp,
		World: "world: " + resp,
	}, common.Ok
}
