package order

import "6-project/internal/product"

type OrderCreateRequest struct {
	ProductsID []uint `json:"products_id" validate:"required,min=1,dive,gt=0"`
}

type OrderResponse struct {
	ID uint `json:"id"`
	UserId uint `json:"userId"`
	Product []product.Product `json:"product"`
}