CREATE TABLE IF NOT EXISTS newsletter_sends (
    id VARCHAR(36) PRIMARY KEY,
    subject TEXT,
    body TEXT,
    status VARCHAR(20),
    sent_count INT DEFAULT 0,
    fail_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);