package auth

type RequestUserByPhone struct {
	Phone string `json:"phone" validate:"required"`
}

type ResponseUserByPhone struct {
	SessionId string `json:"sessionId"`
	Code      int    `json:"code"`
}

type ReqVerifyAuthByCode struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code      int    `json:"code" validate:"required"`
}

type RespTokenAuthByCode struct {
	Token string `json:"token" validate:"required"`
}
