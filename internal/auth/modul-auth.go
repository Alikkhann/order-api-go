package auth

import (
	"gorm.io/gorm"
)

type AuthByPhone struct {
	gorm.Model
	Phone string
	SessionId string
	Code int
}

func NewAuthByPhone(phone, sessionId string, code int) *AuthByPhone {
	return &AuthByPhone{
		Phone: phone,
		SessionId: sessionId,
		Code: code,
	}
}