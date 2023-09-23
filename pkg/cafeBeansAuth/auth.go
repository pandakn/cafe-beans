package cafeBeansAuth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/users"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apiKey"
)

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	// convert nanoseconds to seconds
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

type cafeBeansAuth struct {
	mapClaims *cafeBeansMapClaims
	cfg       config.IJwtConfig
}

type cafeBeansAdmin struct {
	*cafeBeansAuth
}

type cafeBeansApiKey struct {
	*cafeBeansAuth
}

type cafeBeansMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type ICafeBeansAuth interface {
	SignToken() string
}

type ICafeBeansAdmin interface {
	SignToken() string
}

type ICafeBeansApiKey interface {
	SignToken() string
}

func (a *cafeBeansAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *cafeBeansAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func (a *cafeBeansApiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.ApiKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*cafeBeansMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &cafeBeansMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*cafeBeansMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*cafeBeansMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &cafeBeansMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*cafeBeansMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseApiKey(cfg config.IJwtConfig, tokenString string) (*cafeBeansMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &cafeBeansMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*cafeBeansMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &cafeBeansAuth{
		cfg: cfg,
		mapClaims: &cafeBeansMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "cafe-beans-api",
				Subject:   "refresh-token",
				Audience:  []string{"customers", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}

	return obj.SignToken()
}

func NewCafeBeansAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (ICafeBeansAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) ICafeBeansAuth {
	return &cafeBeansAuth{
		cfg: cfg,
		mapClaims: &cafeBeansMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "cafe-beans-api",
				Subject:   "access-token",
				Audience:  []string{"customers", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) ICafeBeansAuth {
	return &cafeBeansAuth{
		cfg: cfg,
		mapClaims: &cafeBeansMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "cafe-beans-api",
				Subject:   "refresh-token",
				Audience:  []string{"customers", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) ICafeBeansAuth {
	return &cafeBeansAdmin{
		cafeBeansAuth: &cafeBeansAuth{
			cfg: cfg,
			mapClaims: &cafeBeansMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "cafe-beans-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300), // 3 minutes
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func newApiKey(cfg config.IJwtConfig) ICafeBeansAuth {
	return &cafeBeansApiKey{
		cafeBeansAuth: &cafeBeansAuth{
			cfg: cfg,
			mapClaims: &cafeBeansMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "cafe-beans-api",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)), // 2 years
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
