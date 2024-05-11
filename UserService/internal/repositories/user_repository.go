package repositories

import (
	"UserSVC/internal/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	CreatePassword(password string) string
	ComparePassword(hashedPwd, password string) error
	FindByID(id string) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreatePassword(password string) string {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}

	return string(hashPwd)
}

func (r *userRepository) ComparePassword(hashedPwd, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
}

func (r *userRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
