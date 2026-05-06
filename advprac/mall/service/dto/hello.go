package dto

type HelloWorldReq struct {
	Name string `json:"name"`
}

type HelloWorldResp struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}
