package appInfoHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/appInfo"
	"github.com/pandakn/cafe-beans/modules/appInfo/appInfoUseCases"
	"github.com/pandakn/cafe-beans/modules/entities"
	"github.com/pandakn/cafe-beans/pkg/cafeBeansAuth"
)

type appInfoHandlerErrCode string

const (
	generateApiKeyErr    appInfoHandlerErrCode = "appInfo-001"
	findCategoryKeyErr   appInfoHandlerErrCode = "appInfo-002"
	insertCategoryKeyErr appInfoHandlerErrCode = "appInfo-003"
	removeCategoryKeyErr appInfoHandlerErrCode = "appInfo-004"
)

type IAppInfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appInfoHandler struct {
	cfg            config.IConfig
	appInfoUseCase appInfoUseCases.IAppInfoUseCase
}

func AppInfoHandler(cfg config.IConfig, appInfoUseCase appInfoUseCases.IAppInfoUseCase) IAppInfoHandler {
	return &appInfoHandler{
		cfg:            cfg,
		appInfoUseCase: appInfoUseCase,
	}
}

func (h *appInfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := cafeBeansAuth.NewCafeBeansAuth(cafeBeansAuth.ApiKey, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appInfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appInfo.CategoryFilter)

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryKeyErr),
			err.Error(),
		).Res()
	}

	categories, err := h.appInfoUseCase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		categories,
	).Res()
}

func (h *appInfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appInfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertCategoryKeyErr),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryKeyErr),
			"categories request are empty",
		).Res()
	}

	if err := h.appInfoUseCase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		req,
	).Res()
}

func (h *appInfoHandler) RemoveCategory(c *fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryKeyErr),
			"id type is invalid",
		).Res()
	}

	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryKeyErr),
			"id must more than 0",
		).Res()
	}

	if err := h.appInfoUseCase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
