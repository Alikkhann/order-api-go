package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Name string	`json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image pq.StringArray `json:"images" validate:"required,dive,url"`
}