package main

import (
	"fmt"
	"log"
	"net/http"
	"6-project/pkg/db"
	"6-project/configs"
	"github.com/gorilla/mux"
	"6-project/internal/auth"
	"6-project/internal/order"
	"6-project/internal/product"
)

func DB() (*db.DB, *configs.Config) {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	return db, conf
}

func App(db *db.DB, conf *configs.Config) http.Handler {
	mux := mux.NewRouter()

	//REPOSITORY
	authRepository := auth.NewAuthRepo(db)
	productRepository := product.NewRepository(db)
	orderRepository := order.NewOrderRepository(db)

	//SERVICE
	serviceAuth := auth.NewServiceAuthByPhone(conf, authRepository)
	serviceOrder := order.NewServiceOrder(&order.OrderServiceDesp{
		OrderRepository:   orderRepository,
		AuthByPhoneRepo:   authRepository,
		ProductRepository: productRepository,
	})

	//HANDLER
	auth.NewHandlerAuthByPhone(mux, &auth.AuthByPhoneHandlerDesp{
		ServiceAuthByPhone: serviceAuth,
		Config:             conf,   //?
	})
	product.NewProductHandler(mux, &product.ProductHandlerDesp{
		ProductRepository: productRepository,
	})

	order.NewHandlerOrder(mux, &order.OrderHandlerDesp{
		Config: conf,
		OrderService:    serviceOrder,
	})

	return mux
}

func main() {
	db, conf := DB()
	app := App(db, conf)
	server := http.Server{
		Addr: ":8081",
		Handler: app,
	}
	fmt.Println("Сервер запущен на порту 8081")
	log.Fatal(server.ListenAndServe())
}