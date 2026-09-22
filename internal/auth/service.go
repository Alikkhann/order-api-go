package auth

import (
	"errors"
	"6-project/configs"
	"6-project/pkg/jwt"
	"6-project/pkg/resp"
)

type ServiceAuthByPhone struct {
	Config *configs.Config 
	Repo *AuthByPhoneRepo
}

func NewServiceAuthByPhone(config *configs.Config, repo *AuthByPhoneRepo) *ServiceAuthByPhone {
	return &ServiceAuthByPhone{
		Config: config,
		Repo: repo,
	}
}

func (repo *ServiceAuthByPhone) AuthByPhone(phone string) (*ResponseUserByPhone, error) {
	checkUser, err := repo.Repo.FindByPhone(phone)
	if checkUser == nil && err == nil {

		sessionId := resp.CreateSessionId()
		code := resp.GenerateVerificationCode()
		user := NewAuthByPhone(phone, sessionId, code)

		err = repo.Repo.Create(user)
		if err != nil {
			return nil, err
		}

		resp := ResponseUserByPhone{SessionId: sessionId, Code: code}
		return &resp, nil
	}

	sessionId := resp.CreateSessionId()
	code := resp.GenerateVerificationCode()

	_, err = repo.Repo.UpdateSessionId(phone, sessionId, code)
	if err != nil {
		return nil, err
	}

	resp := ResponseUserByPhone{SessionId: sessionId, Code: code}
	return &resp, nil
}


func (repo *ServiceAuthByPhone) VerifyByCode(sessionId string, code int) (string, error) {
	user, err := repo.Repo.FindBySessionId(sessionId)
	if err != nil {
		return "", err
	}
	if user == nil || user.Code != code {
		return "", errors.New(ErrInvalidData)
	}
	token, err := jwt.NewJWT(repo.Config.Auth.Secret).Create(&jwt.JWTData{
		Phone: user.Phone,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}