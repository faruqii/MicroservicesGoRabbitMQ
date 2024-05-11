package repositories

import (
	"UserSVC/internal/entities"

	"gorm.io/gorm"
)

type PurchaseRepository interface {
	Create(purchase *entities.Purchase) error
	GetPurchases() ([]entities.Purchase, error)
	GetPurchase(id string) (*entities.Purchase, error)
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *purchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(purchase *entities.Purchase) error {
	return r.db.Create(purchase).Error
}

func (r *purchaseRepository) GetPurchases() ([]entities.Purchase, error) {
	// preload the user and product data
	var purchases []entities.Purchase
	err := r.db.Preload("User").Preload("Product").Find(&purchases).Error
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func (r *purchaseRepository) GetPurchase(id string) (*entities.Purchase, error) {
	var purchase entities.Purchase
	err := r.db.Preload("User").Preload("Product").Where("id = ?", id).First(&purchase).Error
	if err != nil {
		return nil, err
	}

	return &purchase, nil
}
