package service

import (
	"UserSVC/internal/dto"
	"UserSVC/internal/entities"
	"UserSVC/internal/repositories"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	Register(user *entities.User) error
	Login(email, password string) (*entities.User, error)
	CreateUserToken(user *entities.User) (string, error)
	GetAllProducts() ([]entities.Product, error)
	GetProductByID(id string) (*entities.Product, error)
	PurchaseProduct(purchase *entities.Purchase) error
	GetPurchase(id string) (*entities.Purchase, error)
}

type userService struct {
	userRepository     repositories.UserRepository
	productRepository  repositories.ProductRepository
	tokenRepository    repositories.TokenRepository
	purchaseRepository repositories.PurchaseRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
	productRepository repositories.ProductRepository,
	tokenRepository repositories.TokenRepository,
	purchaseRepository repositories.PurchaseRepository,
) *userService {
	return &userService{
		userRepository:     userRepository,
		productRepository:  productRepository,
		tokenRepository:    tokenRepository,
		purchaseRepository: purchaseRepository,
	}
}

func (s *userService) Register(user *entities.User) error {
	// check if email already registered
	_, err := s.userRepository.FindByEmail(user.Email)
	if err == nil {
		return &ErrorMessage{
			Message: "Email already registered",
			Code:    http.StatusInternalServerError,
		}
	}

	// hash password
	user.Password = s.userRepository.CreatePassword(user.Password)

	// create user
	err = s.userRepository.Create(user)
	if err != nil {
		return &ErrorMessage{
			Message: "Failed to create user",
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *userService) Login(email, password string) (*entities.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, &ErrorMessage{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}

	err = s.userRepository.ComparePassword(user.Password, password)
	if err != nil {
		return nil, &ErrorMessage{
			Message: "Invalid password",
			Code:    http.StatusUnauthorized,
		}
	}

	return user, nil
}

func (s *userService) CreateUserToken(user *entities.User) (string, error) {
	// Create JWT token
	claims := dto.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	// Create or update token in repository
	newToken := &entities.Token{
		UserID: user.ID,
		Token:  signedToken,
	}
	_, err = s.tokenRepository.CreateOrUpdateToken(newToken)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *userService) GetAllProducts() ([]entities.Product, error) {
	products, err := s.productRepository.GetAll()
	if err != nil {
		return nil, &ErrorMessage{
			Message: "Failed to get products",
			Code:    http.StatusInternalServerError,
		}
	}

	return products, nil
}

func (s *userService) GetProductByID(id string) (*entities.Product, error) {
	product, err := s.productRepository.GetByID(id)
	if err != nil {
		return nil, &ErrorMessage{
			Message: "Product not found",
			Code:    http.StatusNotFound,
		}
	}

	return product, nil
}

func (s *userService) PurchaseProduct(purchase *entities.Purchase) error {

	// chek if product amount is enough
	err := s.productRepository.CheckProductAmount(purchase.ProductID, purchase.Amount)
	if err != nil {
		return &ErrorMessage{
			Message: "Product amount is not enough",
			Code:    http.StatusBadRequest,
		}
	}

	// check if product exist
	_, err = s.productRepository.GetByID(purchase.ProductID)
	if err != nil {
		return &ErrorMessage{
			Message: "Product not found",
			Code:    http.StatusNotFound,
		}
	}

	// check if user exist
	_, err = s.userRepository.FindByID(purchase.UserID)
	if err != nil {
		return &ErrorMessage{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}

	// create purchase
	err = s.purchaseRepository.Create(purchase)
	if err != nil {
		return &ErrorMessage{
			Message: "Failed to create purchase",
			Code:    http.StatusInternalServerError,
		}
	}

	// update product amount
	err = s.productRepository.UpdateProductAmount(purchase.ProductID, purchase.Amount)
	if err != nil {
		return &ErrorMessage{
			Message: "Failed to update product amount",
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *userService) GetPurchase(id string) (*entities.Purchase, error) {
	purchase, err := s.purchaseRepository.GetPurchase(id)
	if err != nil {
		return nil, &ErrorMessage{
			Message: "Purchase not found",
			Code:    http.StatusNotFound,
		}
	}

	return purchase, nil
}
