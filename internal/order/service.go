package order

import (
	"6-project/internal/auth"
	"6-project/internal/product"
	"errors"
	"gorm.io/gorm"
)

type OrderRepo interface {
    Create(order *Order) error
    GetById(id uint) (*Order, error)
    GetByUserId(userId uint) ([]Order, error)
}

type AuthRepo interface {
    FindByPhone(phone string) (*auth.AuthByPhone, error)
}

type ProductRepo interface {
    CheckManyId(id []uint) ([]product.Product, error)
}

type OrderServiceDesp struct {
	OrderRepository OrderRepo
  AuthByPhoneRepo AuthRepo
	ProductRepository ProductRepo
}

type OrderService struct {
	OrderRepository OrderRepo
  AuthByPhoneRepo AuthRepo
	ProductRepository ProductRepo
}

func NewServiceOrder(desp *OrderServiceDesp) *OrderService {
	return &OrderService{
		OrderRepository:   desp.OrderRepository,
		AuthByPhoneRepo:   desp.AuthByPhoneRepo,
		ProductRepository: desp.ProductRepository,
	}
}

func (repo *OrderService) CreateOrder(body []uint, userData string) (*OrderResponse, error) {
	userId, err := repo.AuthByPhoneRepo.FindByPhone(userData)
	if err != nil {
		return nil, err
	}
	if userId == nil {
		return nil, ErrUserNotFound
	}
	product, err := repo.ProductRepository.CheckManyId(body)
	if err != nil {
		return nil, err
	}
	order := NewOrder(userId.ID, product)
	err = repo.OrderRepository.Create(order)
	if err != nil {
		return nil, err
	}
	resp := OrderResponse{
		ID:      order.ID,
		UserId:  userId.ID,
		Product: product,
	}
	return &resp, nil
}

func (repo *OrderService) GetOrder(phone string, id uint) (*OrderResponse, error) {
	user, err := repo.AuthByPhoneRepo.FindByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	order, err := repo.OrderRepository.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
	return nil, err
	}

	if user.ID != order.UserId {
		return nil, ErrOrderNotFound
	}

	resp := OrderResponse{
		ID: order.ID,
		UserId: order.UserId,
		Product: order.Products,
	}
	return &resp, nil
}

func (repo *OrderService) GetMyOrders(user string) ([]OrderResponse, error)  {
	userId, err := repo.AuthByPhoneRepo.FindByPhone(user)
	if err != nil {
		return nil, err
	}
	if userId == nil {
		return nil, ErrUserNotFound
	}
	orders, err := repo.OrderRepository.GetByUserId(userId.ID)
	if err != nil {
		return nil, err
	}
	var resp []OrderResponse
	for _, order := range orders{
		resp = append(resp, OrderResponse{
			ID: order.ID,
			UserId: order.UserId,
			Product: order.Products,
		})
	}
	return resp, nil
}
