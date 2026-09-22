package auth

import (
	"net/http"
	"6-project/configs"
	"6-project/pkg/req"
	"6-project/pkg/resp"
	"github.com/gorilla/mux"
)

type AuthByPhoneHandlerDesp struct {
	*ServiceAuthByPhone
	*configs.Config
}

type AuthByPhoneHandler struct {
	*ServiceAuthByPhone
}

func NewHandlerAuthByPhone(mux *mux.Router, desp *AuthByPhoneHandlerDesp) {
	handler := AuthByPhoneHandler{
		ServiceAuthByPhone: desp.ServiceAuthByPhone,
	}
	mux.HandleFunc("/authbyphone", handler.Auth())
	mux.HandleFunc("/verifybycode", handler.VerifyCode()).Methods("POST")
}

func (handler *AuthByPhoneHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RequestUserByPhone](w, r)
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return
		}

		user, err := handler.ServiceAuthByPhone.AuthByPhone(body.Phone)
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return
		}

		resp.Json(w, user, 200)
	}
}

func (handler *AuthByPhoneHandler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ReqVerifyAuthByCode](w, r)
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return
		}
		token, err := handler.ServiceAuthByPhone.VerifyByCode(body.SessionId, body.Code)
		if err != nil {
			resp.Json(w, err.Error(), 401)
			return
		}
		data := RespTokenAuthByCode{
			Token: token,
		}
		resp.Json(w, data, 200)
	}
}