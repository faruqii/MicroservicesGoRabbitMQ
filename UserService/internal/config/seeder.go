package config

import (
	"UserSVC/internal/entities"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seed struct {
	DB *gorm.DB
}

func (s *Seed) Seeder() {
	s.UserSeeder()
	s.ProductSeeder()
}

func (s *Seed) UserSeeder() {
	var lenghtTable int64
	s.DB.Model(&entities.User{}).Count(&lenghtTable)
	if lenghtTable == 0 {
		users := []entities.User{
			{
				Name:     "Juanda",
				Email:    "juanda@gmail.com",
				Password: "juanda123",
			},
		}

		for _, user := range users {
			hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
			if err != nil {
				log.Fatalf("Error while hashing password: %v", err)
			}
			user.Password = string(hashPassword)

			err = s.DB.Create(&user).Error
			if err != nil {
				log.Fatalf("Error while seeding user: %v", err)
			}
		}
	}
}

func (s *Seed) ProductSeeder() {
	var lengthTable int64
	s.DB.Model(&entities.Product{}).Count(&lengthTable)
	if lengthTable == 0 {
		products := []entities.Product{
			{
				Name:   "Product 1",
				Price:  100,
				Amount: 10,
			},
			{
				Name:   "Product 2",
				Price:  200,
				Amount: 5,
			},
			// Add more products as needed
		}

		for _, product := range products {
			err := s.DB.Create(&product).Error
			if err != nil {
				log.Fatalf("Failed to seed product: %v", err)
			}
		}
	}
}
