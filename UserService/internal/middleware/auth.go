package middleware

import (
	"UserSVC/internal/repositories"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthenticationMiddleware interface {
	Authenticate() fiber.Handler
}

type Middleware struct {
	tokenRepository repositories.TokenRepository
}

func NewMiddleware(tokenRepository repositories.TokenRepository) *Middleware {
	return &Middleware{
		tokenRepository: tokenRepository,
	}
}

type MiddlewareError struct {
	Message    string
	StatusCode int
}

func (m *Middleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// Extract bearer token from Authorization header
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// Fetch user from token
		user, err := m.tokenRepository.FindUserByToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		c.Locals("user", user)

		// Set bearer token as a custom header
		c.Set("auth", token)

		return c.Next()
	}
}
