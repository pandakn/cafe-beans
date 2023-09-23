package appInfoRepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pandakn/cafe-beans/modules/appInfo"
)

type IAppInfoRepository interface {
	FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error)
	InsertCategory(req []*appInfo.Category) error
	DeleteCategory(categoryId int) error
}

type appInfoRepository struct {
	db *sqlx.DB
}

func AppInfoRepository(db *sqlx.DB) IAppInfoRepository {
	return &appInfoRepository{
		db: db,
	}
}

func (r *appInfoRepository) FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "categories"
	`

	filterValues := make([]any, 0)
	if req.Title != "" {
		query += `
		WHERE (LOWER("title") LIKE $1)
		`

		filterValues = append(filterValues, "%"+strings.ToLower(req.Title)+"%")
	}
	query += ";"

	categories := make([]*appInfo.Category, 0)
	if err := r.db.Select(&categories, query, filterValues...); err != nil {
		return nil, fmt.Errorf("select categories failed: %v", err)
	}

	return categories, nil
}

func (r *appInfoRepository) InsertCategory(req []*appInfo.Category) error {
	ctx := context.Background()

	query := `
	INSERT INTO "categories" (
		"title"
	)
	VALUES
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	valuesStack := make([]any, 0)

	for i, category := range req {
		valuesStack = append(valuesStack, category.Title)

		if i != len(req)-1 {
			query += fmt.Sprintf(`($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`($%d)`, i+1)
		}
	}

	query += `
	RETURNING "id";
	`

	rows, err := tx.QueryxContext(ctx, query, valuesStack...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert categories failed: %v", err)
	}

	var index int
	for rows.Next() {
		if err := rows.Scan(&req[index].Id); err != nil {
			return fmt.Errorf("scan categories id failed %v", err)
		}
		index++
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// func (r *appInfoRepository) DeleteCategory(categoryId int) error {
// 	ctx := context.Background()

// 	query := `DELETE FROM "categories" WHERE "id" = $1;`

// 	if _, err := r.db.ExecContext(ctx, query, categoryId); err != nil {
// 		return fmt.Errorf("delete category failed: %v", err)
// 	}

// 	return nil
// }

func (r *appInfoRepository) DeleteCategory(categoryId int) error {
	ctx := context.Background()

	query := `DELETE FROM "categories" WHERE "id" = $1;`

	result, err := r.db.ExecContext(ctx, query, categoryId)
	if err != nil {
		return fmt.Errorf("delete category failed: %v", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %v", err)
	}

	if rowCount == 0 {
		// No rows were deleted
		return fmt.Errorf("no rows were deleted")
	}

	return nil
}
