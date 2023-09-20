package middlewareRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pandakn/cafe-beans/modules/middleware"
)

type IMiddlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middleware.Role, error)
}

type middlewareRepository struct {
	db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}

func (r *middlewareRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;
	`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}

	return true
}

// Find all role in db
func (r *middlewareRepository) FindRole() ([]*middleware.Role, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;`

	roles := make([]*middleware.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}

	return roles, nil
}
