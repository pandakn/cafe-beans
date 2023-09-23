package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pandakn/cafe-beans/modules/appInfo/appInfoHandlers"
	"github.com/pandakn/cafe-beans/modules/appInfo/appInfoRepositories"
	"github.com/pandakn/cafe-beans/modules/appInfo/appInfoUseCases"
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareHandlers"
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareRepositories"
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareUseCases"
	"github.com/pandakn/cafe-beans/modules/monitor/monitorHandlers"
	"github.com/pandakn/cafe-beans/modules/users/usersHandlers"
	"github.com/pandakn/cafe-beans/modules/users/usersRepositories"
	"github.com/pandakn/cafe-beans/modules/users/usersUseCases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppInfoModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareHandlers.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddleware(s *server) middlewareHandlers.IMiddlewareHandler {
	repository := middlewareRepositories.MiddlewareRepository(s.db)
	useCase := middlewareUseCases.MiddlewareUseCase(repository)
	return middlewareHandlers.MiddlewareHandler(s.cfg, useCase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UserRepository(m.s.db)
	useCase := usersUseCases.UserUseCase(m.s.cfg, repository)
	handler := usersHandlers.UserHandler(m.s.cfg, useCase)

	router := m.r.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SingIn)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)

	// admin
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	// role_id = 2 is admin
	// only admin can access this endpoint
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// user
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppInfoModule() {
	repository := appInfoRepositories.AppInfoRepository(m.s.db)
	useCase := appInfoUseCases.AppInfoUseCase(repository)
	handler := appInfoHandlers.AppInfoHandler(m.s.cfg, useCase)

	router := m.r.Group("/app-info")

	router.Get("/api-key", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	// categories
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}
