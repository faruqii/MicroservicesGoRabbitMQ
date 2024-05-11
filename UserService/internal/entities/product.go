package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID     string `json:"id" gorm:"primaryKey, type:uuid, default:uuid_generate_v4()"`
	Name   string `json:"name"`
	Price  int32  `json:"price"`
	Amount int32  `json:"amount"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}
