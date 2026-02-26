package database

import (
	"chat-app/internal/shared/config"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectSupabaseDB(cfg config.DBConfig) *sqlx.DB {

	db,err:=sqlx.Connect("pgx",	cfg.DSN)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
        log.Fatalf("Database reachable, but ping failed: %v", err)
    }

    log.Println("✅ Database connected & pooling configured")

	return db 

}