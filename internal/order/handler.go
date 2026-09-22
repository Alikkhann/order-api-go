package order

import (
	"6-project/configs"
	"6-project/pkg/contextKeys"
	"6-project/pkg/middleware"
	"6-project/pkg/req"
	"6-project/pkg/resp"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandlerDesp struct {
	*OrderRepository
	*configs.Config
	*OrderService
}

type OrderHandler struct {
	orderService *OrderService
}

func NewHandlerOrder(mux *mux.Router, desp *OrderHandlerDesp) {
	handler := OrderHandler{
		orderService: desp.OrderService,
	}
	mux.Handle("/order", middleware.IsAuthed(handler.Create(), desp.Config)).Methods("POST")
	mux.Handle("/order/{id}", middleware.IsAuthed(handler.Get(), desp.Config)).Methods("GET")
	mux.Handle("/my-orders", middleware.IsAuthed(handler.GetMyOrders(), desp.Config)).Methods("GET")
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[OrderCreateRequest](w, r)
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return
		}
		ctx := r.Context().Value(contextkey.UserContextKey).(string)
		response, err := handler.orderService.CreateOrder(body.ProductsID, ctx)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				resp.Json(w, err.Error(), 404)
				return
			}
			if errors.Is(err, ErrNoProductsFound) {
				resp.Json(w, err.Error(), 400)
				return
			}
			resp.Json(w, err.Error(), 500)
			return
		}

		resp.Json(w, response, 200)
	}
}

func (handler *OrderHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := mux.Vars(r)["id"]
		id, err := strconv.ParseUint(idString, 10, 32) 
		if err != nil {
			resp.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := r.Context().Value(contextkey.UserContextKey).(string)
		order, err := handler.orderService.GetOrder(user, uint(id))
		if err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				resp.Json(w, err.Error(), http.StatusNotFound)
				return
			}
			resp.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Json(w, order, 200)
	}
}

func (handler *OrderHandler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextkey.UserContextKey).(string)
		response, err := handler.orderService.GetMyOrders(user)
		if err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				resp.Json(w, err.Error(), http.StatusNotFound)
			}
			if errors.Is(err, ErrUserNotFound) {
				resp.Json(w, err.Error(), http.StatusNotFound)
			}
			resp.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Json(w, response, 200)
	}
}
