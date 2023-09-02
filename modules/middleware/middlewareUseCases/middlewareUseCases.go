package middlewareUseCases

import (
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareRepositories"
)

type IMiddlewareUseCase interface {
}

type middlewareUseCase struct {
	middlewareRepo middlewareRepositories.IMiddlewareRepository
}

func MiddlewareUseCase(middlewareRepo middlewareRepositories.IMiddlewareRepository) IMiddlewareUseCase {
	return &middlewareUseCase{
		middlewareRepo: middlewareRepo,
	}
}
