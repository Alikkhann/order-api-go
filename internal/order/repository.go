package order

import (
	"6-project/pkg/db"
)

type OrderRepository struct {
	DataBase *db.DB
}

func NewOrderRepository(db *db.DB) *OrderRepository {
	return &OrderRepository{
		DataBase: db,
	}
}

func(repo *OrderRepository) Create(order *Order) error {
	result := repo.DataBase.Create(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func(repo *OrderRepository) GetById(id uint) (*Order, error) {
	var order Order
	result := repo.DataBase.Preload("Products").First(&order, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) GetByUserId(userId uint) ([]Order, error) {
	var orders []Order
	result := repo.DataBase.Preload("Products").Find(&orders, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}