# Digital Wallet API (Go)

A production-grade backend service for a digital wallet application built with Golang. This project demonstrates Clean Architecture, Scalability, and robust Concurrency Control (preventing double-spending) using PostgreSQL transactions and pessimistic locking.

## Key Features

- **Concurrency Safe:** Uses "SELECT ... FOR UPDATE" (Pessimistic Locking) to prevent race conditions during withdrawals.
- **Clean Architecture:** Strict separation of concerns (Handler -> Service -> Repository -> Domain).
- **Scalable:** Database connection pooling and environment-based configuration.
- **Production Ready:** Includes graceful shutdown, structured JSON logging, and standardized error handling.
- **Built-in Migrations:** Custom Go-based migration runner (no external CLI tools required).

## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Router:** Standard Library (net/http)
- **Libraries:**
  - lib/pq (Postgres Driver)
  - joho/godotenv (Environment Variable Management)
  - golang-migrate (Database Migrations)

## Project Structure

wallet/
├── cmd/
│ ├── api/ # Main application entry point
│ └── migrate/ # Database migration utility
├── internal/
│ ├── config/ # Configuration loader
│ ├── domain/ # Business entities & interfaces (Pure Go)
│ ├── handler/ # HTTP Transport layer (Controllers)
│ ├── repository/ # Data Access layer (SQL & Locking)
│ └── service/ # Business Logic layer
├── migrations/ # SQL Migration files (.up.sql / .down.sql)
└── .env # Local environment variables

## Setup and Installation

### 1. Prerequisites

- Go installed
- PostgreSQL installed and running

### 2. Configure Environment

Create a file named ".env" in the root directory with the following content:

DATABASE_URL=postgres://postgres:password@localhost:5432/wallet?sslmode=disable
SERVER_PORT=:8080
DB_MAX_CONN=25

(Note: Replace "postgres", "password", and "wallet" with your actual database credentials.)

### 3. Initialize Database

Run the built-in migration script to create the tables:

go run cmd/migrate/main.go -direction=up

### 4. Run the Server

Start the API:

go run cmd/api/main.go

## API Documentation

### 1. Check Balance

Returns the current balance for a specific user.

- URL: /balance
- Method: GET
- Query Param: user_id (integer)

Example Request:
curl "http://localhost:8080/balance?user_id=1"

Success Response (200 OK):
{
"user_id": 1,
"balance": 1000.00
}

### 2. Withdraw Funds

Safely deducts funds from a user's wallet.

- URL: /withdraw
- Method: POST
- Body: JSON

Example Request:
curl -X POST http://localhost:8080/withdraw \
-H "Content-Type: application/json" \
-d '{"user_id": 1, "amount": 50.00}'

(Windows CMD users: Use double quotes for JSON keys like "{\"user_id\": 1...}")

Success Response (200 OK):
{
"status": "success",
"new_balance": 950.00
}

Error Response (400 Bad Request):
{
"error": "Insufficient funds"
}
