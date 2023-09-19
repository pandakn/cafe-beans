package usersUseCases

import (
	"fmt"

	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/users"
	"github.com/pandakn/cafe-beans/modules/users/usersRepositories"
	"github.com/pandakn/cafe-beans/pkg/cafeBeansAuth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
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

func (u *userUseCase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert a user
	result, err := u.userRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userUseCase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find a user
	user, err := u.userRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is incorrect")
	}

	// Sign Token
	accessToken, _ := cafeBeansAuth.NewCafeBeansAuth(cafeBeansAuth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	refreshToken, _ := cafeBeansAuth.NewCafeBeansAuth(cafeBeansAuth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.userRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *userUseCase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := cafeBeansAuth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check oauth
	oauth, err := u.userRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.userRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := cafeBeansAuth.NewCafeBeansAuth(
		cafeBeansAuth.Access, u.cfg.Jwt(), newClaims,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := cafeBeansAuth.RepeatToken(
		u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.userRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *userUseCase) DeleteOauth(oauthId string) error {
	if err := u.userRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.userRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
