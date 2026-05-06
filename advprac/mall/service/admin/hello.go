package admin

import (
	"context"
	"mall/common"
	"mall/service/dto"
)

func (s *Service) HelloWorld(ctx context.Context, adminUser *common.AdminUser, req *dto.HelloWorldReq) (*dto.HelloWorldResp, common.Errno) {
	return &dto.HelloWorldResp{
		Hello: "hello",
		World: "world",
	}, common.Ok
}
