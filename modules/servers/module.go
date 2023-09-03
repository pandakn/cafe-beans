package servers

import (
	"github.com/gofiber/fiber/v2"
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

	router.Post("/signup", handler.SignUpCustomer)
}
