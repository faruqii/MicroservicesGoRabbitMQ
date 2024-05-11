package routes

import (
	"UserSVC/internal/handlers"
	"UserSVC/internal/middleware"
	"UserSVC/internal/service"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, userSvc service.UserService, mw *middleware.Middleware) {
	userHandler := handlers.NewUserService(userSvc, *mw)

	authRoute := router.Group("/auth")
	authRoute.Post("/register", userHandler.Register)
	authRoute.Post("/login", userHandler.Login)

	productRoute := router.Group("/product")
	productRoute.Get("/", userHandler.GetAllProducts)
	productRoute.Get("/:id", userHandler.GetProductByID)

	purchaseRoute := router.Group("/purchase")
	purchaseRoute.Use(mw.Authenticate())
	purchaseRoute.Post("/", userHandler.PurchaseProduct)
}
