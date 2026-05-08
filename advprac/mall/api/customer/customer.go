package customer

import "mall/adaptor"

type Controller struct {
	adaptor adaptor.IAdaptor

	// services
}

func NewCtrl(adaptor adaptor.IAdaptor) *Controller {
	return &Controller{
		adaptor: adaptor,
	}
}
