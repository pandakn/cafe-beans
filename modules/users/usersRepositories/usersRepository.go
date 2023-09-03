package usersRepositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/pandakn/cafe-beans/modules/users"
	"github.com/pandakn/cafe-beans/modules/users/userPatterns"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
}

type userRepository struct {
	db *sqlx.DB
}

func UserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := userPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get Result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}
