package main

import (
	"io"
	"bytes"
	"testing"
	"net/http"
	"gorm.io/gorm"
	"encoding/json"
	"6-project/pkg/db"
	"github.com/lib/pq"
	"6-project/configs"
	"net/http/httptest"
	"gorm.io/driver/postgres"
	"6-project/internal/auth"
	"6-project/internal/order"
	"6-project/internal/product"
)

func initDB() (*db.DB, error) {
	dsn := "host=localhost user=postgres password=my_pass dbname=orders port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &db.DB{DB: database}, nil
}

func InitData(db *db.DB) uint {
	db.Create(&auth.AuthByPhone{ 
		Phone: "9281112233", 
	})
	product := product.Product{ 
		Name:        "Apple",
		Description: "fruit",
		Images: pq.StringArray{
			"http://example.com/image1.jpg",
			"http://example.com/image2.jpg",
		},
	}
	db.Create(&product)
	return product.ID
}
func RemoveData(db *db.DB, productId uint, orderId uint) {
	db.Exec("Delete FROM order_items WHERE order_id = ?", orderId)
	db.Unscoped().Delete(&product.Product{}, productId)
	db.Unscoped().Where("phone = ?", "9281112233").Delete(&auth.AuthByPhone{})
}

func TestPostOrderSuccess(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatalf("Error create BD: %v", err)
	}
	productID := InitData(db)
	mySecret := "/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="
	conf := &configs.Config{
		Auth: configs.AuthConfig{
			Secret: mySecret, 
		},
	}
	testServ := httptest.NewServer(App(db, conf)) 
	defer testServ.Close()

	login, err := json.Marshal(auth.RequestUserByPhone{
		Phone: "9281112233",
	})
	respLogin, err := http.Post(testServ.URL+"/authbyphone", "application/json", bytes.NewReader(login))
	if err != nil {
		t.Fatalf("login request failed %v", err)
	}
	if respLogin.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", respLogin.StatusCode)
	}
	bodyLogin, err := io.ReadAll(respLogin.Body) 
	if err != nil {
		t.Fatal(err)
	}
	var verifyCode auth.ResponseUserByPhone
	err = json.Unmarshal(bodyLogin, &verifyCode)
	if err != nil {
		t.Fatal(err)
	}

	verifyData, err := json.Marshal(auth.ReqVerifyAuthByCode{
		SessionId: verifyCode.SessionId,
		Code:      verifyCode.Code,
	})
	if err != nil {
		t.Fatal(err)
	}
	respVerify, err := http.Post(testServ.URL+"/verifybycode", "application/json", bytes.NewReader(verifyData))
	if err != nil {
		t.Fatalf("verify request failed %v", err)
	}
	bodyVerify, err := io.ReadAll(respVerify.Body)
	if err != nil {
		t.Fatal(err)
	}
	var Token auth.RespTokenAuthByCode
	err = json.Unmarshal(bodyVerify, &Token)
	if err != nil {
		t.Fatal(err)
	}

	productOrder, err := json.Marshal(order.OrderCreateRequest{
		ProductsID: []uint{productID},
	})
	if err != nil {
		t.Fatal(err)
	}
	reqOrder, err := http.NewRequest("POST", testServ.URL+"/order", bytes.NewReader(productOrder))
	reqOrder.Header.Set("Content-Type", "application/json")
	reqOrder.Header.Set("Authorization", "Bearer "+Token.Token)

	client := &http.Client{}
	respOrder, err := client.Do(reqOrder)
	if err != nil {
		t.Fatal(err)
	}
	if respOrder.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", respOrder.StatusCode)
	}
	bodyOrder, err := io.ReadAll(respOrder.Body)
	if err != nil {
		t.Fatal(err)
	}
	var orderData order.OrderResponse
	err = json.Unmarshal(bodyOrder, &orderData)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		RemoveData(db, productID, orderData.ID)
		sqlDB, _ := db.DB.DB() 
		defer sqlDB.Close() 
	})

	if orderData.ID == 0 {
		t.Error("expect order ID, got 0")
	}
	if len(orderData.Product) == 0 {
		t.Error("expect product in order, got empty")
	}
	if orderData.UserId == 0 {
		t.Error("expect user_id, got 0")
	}
	if len(orderData.Product) > 0 && orderData.Product[0].ID != productID {
		t.Errorf("expect productID %d, got %d", productID, orderData.Product[0].ID)
	}
}
