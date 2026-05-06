package admin

import "mall/adaptor"

type Service struct {
}

func NewService(adaptor *adaptor.Adaptor) *Service {
	return &Service{}
}
