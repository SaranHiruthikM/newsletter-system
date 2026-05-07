package postgres

import (
	"database/sql"
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/domain"
)

type newsletterRepo struct {
	db *sql.DB
}

func NewNewsletterRepository(db *sql.DB) *newsletterRepo {
	return &newsletterRepo{db: db}
}

func (n *newsletterRepo) Create(newsletter *domain.NewsletterSend) error {
	query := `INSERT INTO newsletter_sends (id, subject, body, status, sent_count, fail_count, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := n.db.Exec(query, newsletter.ID, newsletter.Subject, newsletter.Body, newsletter.Status, newsletter.SentCount, newsletter.FailCount, newsletter.CreatedAt, newsletter.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (n *newsletterRepo) IncrementSentCount(id string) error {
	query := `UPDATE newsletter_sends SET sent_count= sent_count+1, updated_at=$2 where id=$1`

	_, err := n.db.Exec(query, id, time.Now())

	if err != nil {
		return err
	}

	return nil

}

func (n *newsletterRepo) IncrementFailCount(id string) error {
	query := `UPDATE newsletter_sends SET fail_count= fail_count+1, updated_at=$2 where id=$1`

	_, err := n.db.Exec(query, id, time.Now())

	if err != nil {
		return err
	}

	return nil

}
