package repositories

import (
	"UserSVC/internal/entities"

	"gorm.io/gorm"
)

type TokenRepository interface {
	CreateOrUpdateToken(token *entities.Token) (string, error)
	GetTokenByUserID(userID string) (*entities.Token, error)
	FindUserByToken(token string) (string, error)
	GetUserIDByToken(token string) (string, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) CreateOrUpdateToken(token *entities.Token) (string, error) {
	var existingToken entities.Token
	err := r.db.Where("user_id = ?", token.UserID).First(&existingToken).Error

	if err != nil {
		if err := r.db.Create(token).Error; err != nil {
			return "", err
		}
	} else {
		token.ID = existingToken.ID
		if err := r.db.Save(token).Error; err != nil {
			return "", err
		}
	}

	return token.Token, nil
}

func (r *tokenRepository) GetTokenByUserID(userID string) (*entities.Token, error) {
	var token entities.Token
	if err := r.db.Where("user_id = ?", userID).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) FindUserByToken(token string) (string, error) {
	var user entities.Token
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.UserID, nil
}

func (r *tokenRepository) GetUserIDByToken(token string) (string, error) {
	var user entities.Token
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.UserID, nil
}
