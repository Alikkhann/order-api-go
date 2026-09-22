package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Phone string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data *JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ 
		"phone": data.Phone,
	})
	signed, err := token.SignedString([]byte(j.Secret)) 
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {	
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	phone := tkn.Claims.(jwt.MapClaims)["phone"] 
	return tkn.Valid, &JWTData{
		Phone: phone.(string),
	}
}