package product

import (
	"6-project/pkg/db"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	DataBase *db.DB
}

func NewRepository(db *db.DB) *ProductRepository{
	return &ProductRepository{
		DataBase: db,
	}
}

func (create *ProductRepository) Create(product *Product) (*Product, error){
	result := create.DataBase.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (get *ProductRepository) GetById(id uint) (*Product, error) {
	product := Product{}
	result := get.DataBase.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (update *ProductRepository) Update(product *Product) (*Product, error) {
	result := update.DataBase.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (delete *ProductRepository) Delete(id uint) error {
	result := delete.DataBase.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (check *ProductRepository) CheckId(id uint) error {
	var product Product
	result := check.DataBase.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (check *ProductRepository) CheckManyId(id []uint) ([]Product, error) {
	var products []Product
	result := check.DataBase.DB.Find(&products, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}