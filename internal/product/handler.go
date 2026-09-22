package product

import (
	"strconv"
	"net/http"
	"gorm.io/gorm"
	"6-project/pkg/req"
	"6-project/pkg/resp"
	"github.com/gorilla/mux"
)

type ProductHandlerDesp struct {
	ProductRepository *ProductRepository
}

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(mux *mux.Router, desp *ProductHandlerDesp) {
	handler := ProductHandler{
		ProductRepository: desp.ProductRepository,
	}
	mux.HandleFunc("/product", handler.Create()).Methods("POST")
	mux.HandleFunc("/product/{id}", handler.GoTo()).Methods("GET")
	mux.HandleFunc("/product/{id}", handler.Update()).Methods("PATCH")
	mux.HandleFunc("/product/{id}", handler.Delete()).Methods("DELETE")
}

func (handler *ProductHandler) Create() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		productBody, err := req.HandleBody[ProductCreateRequest](w, r) 
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return 
		}
		product := NewProduct(*productBody)
		createProduct, err := handler.ProductRepository.Create(product)
		if err != nil {
			resp.Json(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		resp.Json(w, createProduct, 201)
	}
}

func (handler *ProductHandler) GoTo() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			resp.Json(w, "Invalid product ID", http.StatusBadRequest)
			return 
		}
		product, err := handler.ProductRepository.GetById(uint(id))
		if err != nil {
			resp.Json(w, err.Error(), http.StatusNotFound)
			return
		}
		resp.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		productBody, err := req.HandleBody[ProductCreateRequest](w, r)
		if err != nil {
			resp.Json(w, err.Error(), 400)
			return 
		}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			resp.Json(w, "Invalid product ID", http.StatusBadRequest)
			return 
		}
		product, err := handler.ProductRepository.Update(&Product{
			Model: gorm.Model{ID: uint(id)},
			Name: productBody.Name,
			Description: productBody.Description,
			Images: productBody.Image,
		})
		if err != nil {
			resp.Json(w, err.Error(), http.StatusBadRequest)
			return 
		}
		resp.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			resp.Json(w, "Invalid product ID", http.StatusBadRequest)
			return 
		}
		err = handler.ProductRepository.CheckId(uint(id))
		if err != nil {
			resp.Json(w, err.Error(), http.StatusNotFound)
			return 
		}
		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			resp.Json(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		resp.Json(w, nil, 200)
	}
}