package dto

import "github.com/golang-jwt/jwt/v4"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type PurchaseProductRequest struct {
	ProductID string `json:"product_id"`
	Amount    int32  `json:"amount"`
}

type PurchaseProductResponse struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Product string `json:"product"`
	Amount  int32  `json:"amount"`
}
