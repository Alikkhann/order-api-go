package product

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Product struct {
  gorm.Model
  Name        string
  Description string
  Images      pq.StringArray `gorm:"type:text[]"`
}

func NewProduct(prodBody ProductCreateRequest) *Product{
  return &Product{
    Name: prodBody.Name,
    Description: prodBody.Description,
    Images: prodBody.Image,
  }
}


