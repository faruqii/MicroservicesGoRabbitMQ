package repositories

import (
	"UserSVC/internal/entities"
	"errors"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]entities.Product, error)
	GetByID(id string) (*entities.Product, error)
	UpdateProductAmount(id string, amount int32) error
	// check if product amount is enough
	CheckProductAmount(id string, amount int32) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetByID(id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) UpdateProductAmount(id string, purchasedAmount int32) error {
	var product entities.Product
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if there are enough products in stock
	if product.Amount < purchasedAmount {
		return errors.New("not enough stock available")
	}

	// Decrease the stock
	product.Amount -= purchasedAmount

	// Update the product amount in the database
	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) CheckProductAmount(id string, amount int32) error {
	var product entities.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return err
	}

	if product.Amount < amount {
		return err
	}

	return nil
}
