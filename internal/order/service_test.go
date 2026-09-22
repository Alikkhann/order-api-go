package order_test

import (
	"6-project/internal/auth"
	"6-project/internal/order"
	"6-project/internal/product"
	"errors"
	"testing"
	"gorm.io/gorm"
)

type MokAuthRepo struct {
	ReturnNil bool
}
type MokProductRepo struct {
	ReturnError bool
}
type MokOrderRepo struct {
	ReturnError bool
}

func (repo *MokAuthRepo) FindByPhone(phone string) (*auth.AuthByPhone, error) {
	if repo.ReturnNil == true {
		return nil, nil
	}
	return &auth.AuthByPhone{
		Model: gorm.Model{ID: 1},
		Phone: phone,
	}, nil
}

func (repo *MokProductRepo) CheckManyId(id []uint) ([]product.Product, error) {
	if repo.ReturnError == true {
		return nil, errors.New("product error")
	}
	return []product.Product{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "Apple",
			Description: "Fruit",
		},
	}, nil
}

func (repo *MokOrderRepo) Create(order *order.Order) error {
	if repo.ReturnError == true {
		return errors.New("create error")
	}
	return nil
}

func (repo *MokOrderRepo) GetById(id uint) (*order.Order, error) {
	if repo.ReturnError == true {
		return nil, gorm.ErrRecordNotFound
	}
	return &order.Order{
		Model:  gorm.Model{ID: 1},
		UserId: 1,
		Products: []product.Product{
			{
				Model: gorm.Model{ID: id},
				Name:  "Apple",
			},
		},
	}, nil
}

func (repo *MokOrderRepo) GetByUserId(userId uint) ([]order.Order, error) {
	if repo.ReturnError == true {
		return nil, errors.New("user not found")
	}
	return []order.Order{
		{
			Model:  gorm.Model{ID: 1},
			UserId: 1,
			Products: []product.Product{
				{
					Model: gorm.Model{ID: 1},
					Name:  "Apple",
				},
			},
		},
	}, nil
}

func NewTestService() *order.OrderService {
	return &order.OrderService{
		OrderRepository:   &MokOrderRepo{},
		AuthByPhoneRepo:   &MokAuthRepo{},
		ProductRepository: &MokProductRepo{},
	}
}

func TestServiceCreateOrder_Success(t *testing.T) {
	service := NewTestService()
	resp, err := service.CreateOrder([]uint{1}, "9281112233")
	if err != nil {
		t.Fatal(err)
	}
	if resp.UserId != 1 {
		t.Fatalf("expected userId 1, got %d", resp.UserId)
	}
	if len(resp.Product) == 0 {
		t.Errorf("expected products, got empty")
	}
}


func TestServiceGetOrder_Success(t *testing.T) {
	service := NewTestService()
	resp, err := service.GetOrder("9281112233", 1)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Errorf("expected order, get empty")
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
	if len(resp.Product) == 0 {
		t.Error("expected products, got empty")
	}
}

func TestServiceGetMyOrders_Success(t *testing.T) {
	service := NewTestService()
	resp, err := service.GetMyOrders("9281112233")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("expected orders, got empty")
	}
	if resp[0].UserId != 1 {
		t.Errorf("expected userId 1, got %d", resp[0].UserId)
	}
	if len(resp[0].Product) == 0 {
		t.Error("expected products, got empty")
	}
}

func TestServiceCreateOrder_UserNotFound(t *testing.T) {
    service := &order.OrderService{
        OrderRepository:   &MokOrderRepo{},
        AuthByPhoneRepo:   &MokAuthRepo{ReturnNil: true},
        ProductRepository: &MokProductRepo{},
    }

    _, err := service.CreateOrder([]uint{1}, "9281112233")
    if !errors.Is(err, order.ErrUserNotFound) {
        t.Errorf("expected ErrUserNotFound, got %v", err)
    }
}


func TestServiceCreateOrder_ProductError(t *testing.T) {
    service := &order.OrderService{
        OrderRepository:   &MokOrderRepo{},
        AuthByPhoneRepo:   &MokAuthRepo{},
        ProductRepository: &MokProductRepo{ReturnError: true},
    }

    _, err := service.CreateOrder([]uint{1}, "9281112233")
    if err == nil {
        t.Error("expected error, got nil")
    }
}


func TestServiceCreateOrder_CreateError(t *testing.T) {
    service := &order.OrderService{
        OrderRepository:   &MokOrderRepo{ReturnError: true},
        AuthByPhoneRepo:   &MokAuthRepo{},
        ProductRepository: &MokProductRepo{},
    }

    _, err := service.CreateOrder([]uint{1}, "9281112233")
    if err == nil {
        t.Error("expected error, got nil")
    }
}

func TestServiceGetOrder_NotFound(t *testing.T) {
	service := &order.OrderService{
    OrderRepository:   &MokOrderRepo{ReturnError: true},
    AuthByPhoneRepo:   &MokAuthRepo{},
    ProductRepository: &MokProductRepo{},
  }
	_, err := service.GetOrder("9281112233", uint(1))
	if !errors.Is(err, order.ErrOrderNotFound) {
		t.Errorf("expected ErrorNotFound, got %v", err)
	}
}

func TestServiceGetMyOrders_UserNotFound(t *testing.T) {
	service := &order.OrderService{
    OrderRepository:   &MokOrderRepo{},
    AuthByPhoneRepo:   &MokAuthRepo{ReturnNil: true},
    ProductRepository: &MokProductRepo{},
  }
	_, err := service.GetMyOrders("9281112233")
	if !errors.Is(err, order.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}