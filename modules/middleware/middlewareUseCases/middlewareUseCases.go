package middlewareUseCases

import (
	"github.com/pandakn/cafe-beans/modules/middleware"
	"github.com/pandakn/cafe-beans/modules/middleware/middlewareRepositories"
)

type IMiddlewareUseCase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middleware.Role, error)
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

func (u *middlewareUseCase) FindRole() ([]*middleware.Role, error) {
	roles, err := u.middlewareRepo.FindRole()

	if err != nil {
		return nil, err
	}

	return roles, nil
}
