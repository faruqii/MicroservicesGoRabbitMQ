package app

import (
	"UserSVC/internal/config"
	"UserSVC/internal/middleware"
	"UserSVC/internal/repositories"
	"UserSVC/internal/routes"
	"UserSVC/internal/service"
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartService() {
	app := fiber.New()

	db, err := config.Connect()
	if err != nil {
		panic(err)
	}

	seed := config.Seed{DB: db}
	seed.Seeder()

	// init repo
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)

	// init service
	userSvc := service.NewUserService(userRepo, productRepo, tokenRepo, purchaseRepo)

	// middleware
	middleware := middleware.NewMiddleware(tokenRepo)

	// route
	api := app.Group("/api")
	routes.UserRoutes(api, userSvc, middleware)

	err = app.Listen(":" + os.Getenv("PORT"))

	if err != nil {
		panic(err)
	}
}
