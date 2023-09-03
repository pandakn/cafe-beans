package usersUseCases

import (
	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/users"
	"github.com/pandakn/cafe-beans/modules/users/usersRepositories"
)

type IUserUseCase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type userUseCase struct {
	cfg            config.IConfig
	userRepository usersRepositories.IUserRepository
}

func UserUseCase(cfg config.IConfig, userRepository usersRepositories.IUserRepository) IUserUseCase {
	return &userUseCase{
		cfg:            cfg,
		userRepository: userRepository,
	}
}
func (u *userUseCase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert a user
	result, err := u.userRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
