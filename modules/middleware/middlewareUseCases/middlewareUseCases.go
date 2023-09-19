package middlewareUseCases

import (
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareRepositories"
)

type IMiddlewareUseCase interface {
	FindAccessToken(userId, accessToken string) bool
}

type middlewareUseCase struct {
	middlewareRepo middlewareRepositories.IMiddlewareRepository
}

func MiddlewareUseCase(middlewareRepo middlewareRepositories.IMiddlewareRepository) IMiddlewareUseCase {
	return &middlewareUseCase{
		middlewareRepo: middlewareRepo,
	}
}

func (u *middlewareUseCase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepo.FindAccessToken(userId, accessToken)
}
