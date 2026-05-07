package database

import (
	"database/sql"
	"fmt"
	"github.com/SaranHiruthikM/newsletter-system/internal/config"
	_ "github.com/lib/pq"
)

func Connect(cfg config.DBConfig) (*sql.DB, error) {
	conn_string := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	conn, err := sql.Open("postgres", conn_string)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
