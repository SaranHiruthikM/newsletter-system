package postgres

import "database/sql"

type newsletterRepo struct {
	db *sql.DB
}

func NewNewsletterRepository(db *sql.DB) *subscriberRepo {
	return &subscriberRepo{db: db}
}
