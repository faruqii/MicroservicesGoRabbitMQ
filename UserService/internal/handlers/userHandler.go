package handlers

import (
	"UserSVC/internal/dto"
	"UserSVC/internal/entities"
	"UserSVC/internal/middleware"
	"UserSVC/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type UserHandlers struct {
	userSvc           service.UserService
	middlewareManager middleware.Middleware
}

func NewUserService(userSvc service.UserService, middlewareManager middleware.Middleware) *UserHandlers {
	return &UserHandlers{
		userSvc:           userSvc,
		middlewareManager: middlewareManager,
	}
}

func (h *UserHandlers) Register(ctx *fiber.Ctx) (err error) {
	var req dto.RegisterRequest
	if err = ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := &entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err = h.userSvc.Register(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func (h *UserHandlers) Login(ctx *fiber.Ctx) (err error) {
	var req dto.LoginRequest
	if err = ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userSvc.Login(req.Email, req.Password)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := h.userSvc.CreateUserToken(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login success",
		"data":    response,
		"token":   token,
	})
}

func (h *UserHandlers) GetAllProducts(ctx *fiber.Ctx) (err error) {
	products, err := h.userSvc.GetAllProducts()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    products,
	})
}

func (h *UserHandlers) GetProductByID(ctx *fiber.Ctx) (err error) {
	id := ctx.Params("id")
	product, err := h.userSvc.GetProductByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    product,
	})
}

func (h *UserHandlers) PurchaseProduct(ctx *fiber.Ctx) (err error) {
	user := ctx.Locals("user").(string)

	var req dto.PurchaseProductRequest
	if err = ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	purchase := &entities.Purchase{
		UserID:    user,
		ProductID: req.ProductID,
		Amount:    req.Amount,
	}

	err = h.userSvc.PurchaseProduct(purchase)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retrieve the complete purchase data with product name
	completePurchase, err := h.userSvc.GetPurchase(purchase.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := dto.PurchaseProductResponse{
		ID:      completePurchase.ID,
		User:    user,
		Product: completePurchase.Product.Name, // Access product name from the complete purchase data
		Amount:  completePurchase.Amount,
	}

	// Publish message to RabbitMQ after a successful purchase
	purchaseJSON, err := json.Marshal(response)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Publish the purchase data to RabbitMQ exchange
	err = h.PublishToRabbitMQ(purchaseJSON)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Purchase success",
		"data":    response,
	})
}

func (h *UserHandlers) PublishToRabbitMQ(data []byte) error {
	// Initialize RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"pubsub", // queue name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	body := data
	err = channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
