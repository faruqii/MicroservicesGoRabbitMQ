package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Purchase struct {
	ID        string  `json:"id" gorm:"primaryKey, type:uuid, default:uuid_generate_v4()"`
	UserID    string  `json:"user_id"`
	User      User    `json:"user"`
	ProductID string  `json:"product_id"`
	Product   Product `json:"product"`
	Amount    int32   `json:"amount"`
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}
