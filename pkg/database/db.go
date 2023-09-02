package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pandakn/cafe-beans/config"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", cfg.Url()) // require import pgx cuz sqlx call driver pgx
	if err != nil {
		log.Fatalf("connect to db failed %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())

	return db
}
