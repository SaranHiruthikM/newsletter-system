# Newsletter System

This project helps you manage email subscriptions and send out newsletters to your audience. It takes care of new sign-ups, sends confirmation emails, and lets you broadcast messages to all your confirmed subscribers, with everything running efficiently in the background.

## Description

Ever needed a reliable way to connect with your subscribers? This system provides a robust and scalable solution for managing email newsletters. Built with performance in mind, it handles everything from user subscription and email confirmation to the asynchronous dispatch of bulk newsletters. With features like rate limiting and idempotency, it ensures a smooth and reliable experience for both you and your subscribers.

## Installation

To get this project up and running locally, follow these steps:

1.  **Clone the Repository**:

    ```bash
    git clone https://github.com/SaranHiruthikM/newsletter-system.git
    cd newsletter-system
    ```

2.  **Prerequisites**:
    - Go (version 1.25.9 or higher)
    - PostgreSQL
    - Redis
    - RabbitMQ

    Ensure that PostgreSQL, Redis, and RabbitMQ are running and accessible.

3.  **Environment Configuration**:
    Copy the example environment file and fill in your details:

    ```bash
    cp .env.example .env
    ```

    Edit the `.env` file with your database, Redis, RabbitMQ, Resend API key, and other application settings.

    ```ini
    APP_PORT=3001
    APP_ENV=development

    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=newsletter_db
    DB_SSLMODE=disable

    REDIS_HOST=localhost
    REDIS_PORT=6379

    RABBITMQ_URL=amqp://guest:guest@localhost:5672/ # Example for local RabbitMQ

    RATE_LIMIT_ENABLED=true
    RATE_LIMIT_MAX_REQUESTS=5
    RATE_LIMIT_WINDOW=1m

    IDEMPOTENCY_ENABLED=true
    IDEMPOTENCY_TTL=10m

    EMAIL_PROVIDER=resend
    EMAIL_FROM_EMAIL=news@example.com
    EMAIL_FROM_NAME=Your Newsletter
    RESEND_API_KEY=your_resend_api_key
    RESEND_BASE_URL=https://api.resend.com
    RESEND_TIMEOUT=10s

    ADMIN_API_KEY=your_admin_secret_key
    WORKER_METRICS_PORT=3002
    ```

4.  **Database Setup**:
    Ensure your PostgreSQL database is created and running. You'll need to manually create the necessary tables for subscribers and newsletters.
    Here's an example schema:

    ```sql
    CREATE TABLE subscribers (
        id VARCHAR(255) PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        confirmed BOOLEAN NOT NULL DEFAULT FALSE,
        token VARCHAR(255) NOT NULL,
        token_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL
    );

    CREATE TABLE newsletter_sends (
        id VARCHAR(255) PRIMARY KEY,
        subject VARCHAR(255) NOT NULL,
        body TEXT NOT NULL,
        status VARCHAR(50) NOT NULL,
        sent_count INTEGER NOT NULL DEFAULT 0,
        fail_count INTEGER NOT NULL DEFAULT 0,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    ```

5.  **Install Dependencies**:

    ```bash
    go mod tidy
    ```

6.  **Run the API Server**:
    The API listens on the `APP_PORT` defined in your `.env` file (default `3001`).

    ```bash
    go run cmd/api/main.go
    ```

7.  **Run the Worker**:
    The worker processes email sending asynchronously and listens on `WORKER_METRICS_PORT` (default `3002`) for metrics.
    ```bash
    go run cmd/worker/main.go
    ```

## Usage

Once both the API and worker are running, you can interact with the system:

1.  **Subscribe a User**:
    Send a POST request to `/api/v1/subscribe` with the user's email. A confirmation email will be queued and sent by the worker.

    **Request**:
    `POST /api/v1/subscribe`

    ```json
    {
      "email": "test@example.com"
    }
    ```

    **Response (201 Created)**:

    ```json
    {
      "message": "subscription successful, please check your mail"
    }
    ```

2.  **Confirm Subscription**:
    The user will receive an email with a confirmation link. This link typically looks like `http://localhost:8080/api/v1/confirm?token=YOUR_TOKEN`.
    When the user clicks the link, a GET request is made to `/api/v1/confirm` with the `token` as a query parameter.

    **Request**:
    `GET /api/v1/confirm?token=some-unique-token-from-email`

    **Response (200 OK)**:

    ```json
    {
      "message": "email confimred successfully"
    }
    ```

3.  **Send a Newsletter**:
    To send a newsletter, you'll need your `ADMIN_API_KEY` in the `X-API-Key` header. This will queue a newsletter to all _confirmed_ subscribers.

    **Request**:
    `POST /api/v1/newsletter/send`
    `X-API-Key: your_admin_secret_key` (from `.env`)
    `Idempotency-Key: a-unique-key-for-this-request` (required if idempotency is enabled)

    ```json
    {
      "subject": "Exciting News!",
      "body": "<html><body><h1>Hello Subscriber!</h1><p>Check out our latest updates.</p></body></html>"
    }
    ```

    **Response (200 OK)**:

    ```json
    {
      "message": "newsletter dispatch started",
      "total": 500
    }
    ```

4.  **Monitoring with Prometheus and Grafana**:
    - The API exposes metrics at `/api/v1/metrics`.
    - The worker exposes metrics at `/metrics` on `WORKER_METRICS_PORT`.
    - You can set up Prometheus to scrape these endpoints using the provided `deploy/prometheus.yml` configuration.
    - Import the `deploy/grafana-dashboards/API scrape health-1778323970213.json` file into Grafana to visualize key metrics like email send rates, failures, and processing times.

## Features

- **Subscriber Management**: Easily handle new subscriptions, confirmation flows, and maintain a list of confirmed users.
- **Asynchronous Email Delivery**: Utilizes RabbitMQ to queue email sending tasks, ensuring the API remains responsive while emails are processed by a dedicated worker.
- **Email Confirmation**: Implements a token-based email confirmation system to verify subscriber authenticity.
- **Newsletter Broadcasting**: Send newsletters to all confirmed subscribers efficiently.
- **Rate Limiting**: Protects public API endpoints from abuse using Redis-backed rate limiting.
- **Idempotency**: Prevents duplicate processing of critical requests (e.g., sending newsletters) using Redis-backed idempotency keys.
- **API Key Authentication**: Secures administrative endpoints with a simple API key.
- **Health Checks**: Provides a `/api/v1/health` endpoint to monitor database connectivity.
- **Observability**: Integrated Prometheus metrics for real-time monitoring and a Grafana dashboard for visualization of worker performance.
- **External Email Provider Integration**: Supports email sending via Resend.com, extensible for other providers.

## Technologies Used

| Technology     | Description                                                       |
| :------------- | :---------------------------------------------------------------- |
| **Go**         | Primary language for backend development.                         |
| **Fiber**      | Fast, expressive web framework for Go.                            |
| **PostgreSQL** | Robust relational database for data persistence.                  |
| **Redis**      | In-memory data store for caching, rate limiting, and idempotency. |
| **RabbitMQ**   | Message broker for asynchronous task processing.                  |
| **Prometheus** | Monitoring system for collecting metrics.                         |
| **Grafana**    | Data visualization and dashboarding tool.                         |
| **Resend**     | Transactional email API for reliable email delivery.              |

## Contributing

We welcome contributions! If you have suggestions for improvements, new features, or bug fixes, please feel free to:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and ensure tests pass.
4.  Commit your changes with clear, descriptive messages.
5.  Push your branch to your forked repository.
6.  Open a pull request to the `main` branch of this repository.

Please make sure to follow the existing code style and conventions.
