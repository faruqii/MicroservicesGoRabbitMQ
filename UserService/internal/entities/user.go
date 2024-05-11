package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey, type:uuid, default:uuid_generate_v4()"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}
