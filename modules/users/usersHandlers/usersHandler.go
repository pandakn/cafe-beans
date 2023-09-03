package usersHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/entities"
	"github.com/pandakn/cafe-beans/modules/users"
	"github.com/pandakn/cafe-beans/modules/users/usersUseCases"
)

type usersHandlersErrCode string

const (
	signUpCustomerErr usersHandlersErrCode = "users-001"
)

type IUserHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type userHandler struct {
	cfg         config.IConfig
	userUseCase usersUseCases.IUserUseCase
}

func UserHandler(cfg config.IConfig, userUseCase usersUseCases.IUserUseCase) IUserHandler {
	return &userHandler{
		cfg:         cfg,
		userUseCase: userUseCase,
	}
}

func (h *userHandler) SignUpCustomer(c *fiber.Ctx) error {
	// request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// validate email
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email is not a valid email address",
		).Res()
	}

	// Insert
	result, err := h.userUseCase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()

		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
