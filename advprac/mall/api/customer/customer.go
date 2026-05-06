package customer

import "mall/adaptor"

type Controller struct {
	adaptor *adaptor.Adaptor

	// services
}

func NewCtrl(adaptor *adaptor.Adaptor) *Controller {
	return &Controller{
		adaptor: adaptor,
	}
}
