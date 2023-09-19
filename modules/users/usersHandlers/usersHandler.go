package usersHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/entities"
	"github.com/pandakn/cafe-beans/modules/users"
	"github.com/pandakn/cafe-beans/modules/users/usersUseCases"
	"github.com/pandakn/cafe-beans/pkg/cafeBeansAuth"
)

type usersHandlersErrCode string

const (
	signUpCustomerErr  usersHandlersErrCode = "users-001"
	signInErr          usersHandlersErrCode = "users-002"
	refreshPassportErr usersHandlersErrCode = "users-003"
	signOutErr         usersHandlersErrCode = "users-004"
	signUpAdminErr     usersHandlersErrCode = "users-005"
	generateAdminErr   usersHandlersErrCode = "users-006"
	getUserProfileErr  usersHandlersErrCode = "users-007"
)

type IUserHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	SingIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
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

func (h *userHandler) SignUpAdmin(c *fiber.Ctx) error {
	// request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
	}

	// validate email
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			"email is not a valid email address",
		).Res()
	}

	// Insert
	result, err := h.userUseCase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpAdminErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpAdminErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpAdminErr),
				err.Error(),
			).Res()

		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *userHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := cafeBeansAuth.NewCafeBeansAuth(cafeBeansAuth.Admin,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateAdminErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *userHandler) SingIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUseCase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUseCase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.userUseCase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "Logout successfully").Res()
}

func (h *userHandler) GetUserProfile(c *fiber.Ctx) error {
	// Set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	result, err := h.userUseCase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
