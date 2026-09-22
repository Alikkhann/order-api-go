package auth

import (
	"errors"
	"gorm.io/gorm"
	"6-project/pkg/db"
)

type AuthByPhoneRepo struct {
	AuthRepo *db.DB
}

func NewAuthRepo(db *db.DB) *AuthByPhoneRepo {
	return &AuthByPhoneRepo{
		AuthRepo: db,
	}
}

func (repo *AuthByPhoneRepo) Create(user *AuthByPhone) error {
	result := repo.AuthRepo.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *AuthByPhoneRepo) FindByPhone(number string) (*AuthByPhone, error) {
	var user AuthByPhone
	result := repo.AuthRepo.DB.First(&user, "phone = ?", number)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *AuthByPhoneRepo) UpdateSessionId(number, sessionId string, code int) (*AuthByPhone, error) {
	var user AuthByPhone
	result := repo.AuthRepo.DB.First(&user, "phone = ?", number)
	if result.Error != nil {
		return nil, result.Error
	}
	user.SessionId = sessionId
	user.Code = code
	resultUp := repo.AuthRepo.DB.Model(&user).Updates(&user)
	if resultUp.Error != nil {
		return nil, resultUp.Error
	}
	return &user, nil
}

func (repo *AuthByPhoneRepo) FindBySessionId(sessionId string) (*AuthByPhone, error) {
	var user AuthByPhone
	result := repo.AuthRepo.DB.First(&user, "session_Id = ?", sessionId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
