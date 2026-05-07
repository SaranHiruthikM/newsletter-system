package postgres

import (
	"database/sql"
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/domain"
)

type subscriberRepo struct {
	db *sql.DB
}

func NewSubscriberRepository(db *sql.DB) *subscriberRepo {
	return &subscriberRepo{db: db}
}

func (s *subscriberRepo) FindByEmail(email string) (*domain.Subscriber, error) {
	query := `SELECT id, email, confirmed, token, token_expires_at, created_at, updated_at 
              FROM subscribers WHERE email = $1`

	row := s.db.QueryRow(query, email)

	sub := &domain.Subscriber{}

	err := row.Scan(
		&sub.ID,
		&sub.Email,
		&sub.Confirmed,
		&sub.Token,
		&sub.TokenExpiresAt,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriberRepo) Create(subscriber *domain.Subscriber) error {
	query := `INSERT INTO subscribers (id, email, confirmed, token, token_expires_at, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(query, subscriber.ID, subscriber.Email, subscriber.Confirmed, subscriber.Token, subscriber.TokenExpiresAt, subscriber.CreatedAt, subscriber.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *subscriberRepo) FindByToken(token string) (*domain.Subscriber, error) {

	query := `SELECT id, email, confirmed, token, token_expires_at, created_at, updated_at 
              FROM subscribers WHERE token = $1`

	row := s.db.QueryRow(query, token)

	sub := &domain.Subscriber{}

	err := row.Scan(
		&sub.ID,
		&sub.Email,
		&sub.Confirmed,
		&sub.Token,
		&sub.TokenExpiresAt,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriberRepo) UpdateConfirmed(id string, confirmed bool) error {

	query := `UPDATE subscribers SET confirmed=$2, updated_at=$3 WHERE id=$1`

	_, err := s.db.Exec(query, id, confirmed, time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (s *subscriberRepo) FindAllConfirmed() ([]*domain.Subscriber, error) {
	var subs []*domain.Subscriber

	query := `SELECT id, email, confirmed, token, token_expires_at, created_at, updated_at 
              FROM subscribers WHERE confirmed=true `

	rows, err := s.db.Query(query)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		sub := &domain.Subscriber{}

		if err := rows.Scan(
			&sub.ID,
			&sub.Email,
			&sub.Confirmed,
			&sub.Token,
			&sub.TokenExpiresAt,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		); err != nil {
			return nil, err
		}

		subs = append(subs, sub)

	}

	return subs, nil
}
