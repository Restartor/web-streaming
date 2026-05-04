# Web Streaming Backend

This repository contains a Go backend for a simple web streaming app. It uses Gin for HTTP routing, GORM with PostgreSQL for persistence, JWT for authentication, and rate limiting for selected public endpoints.

## Features

- User registration and login
- JWT-based authentication
- Admin-only film management
- Film listing and title search
- PostgreSQL-backed storage with automatic migration on startup

## Requirements

- Go 1.25 or newer
- PostgreSQL database
- A `.env` file with the required environment variables

## Setup

1. Open a terminal in the `backend` folder.

2. Install dependencies.

	```bash
	go mod tidy
	```

3. Create a `.env` file in `backend` with your database and JWT settings.

4. Start the server.

	```bash
	go run .
	```

The server starts on port `1010` by default.

## Environment Variables

The application reads these values at startup:

- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_PORT`
- `DB_SSLMODE`
- `JWT_SECRET`

Example:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=web_streaming
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
```

## API Routes

Public routes under `/api/v1`:

- `POST /register` - create a new user
- `POST /login` - sign in and receive a JWT token
- `GET /films` - list all films
- `GET /films/search?title=...` - search films by title

Protected admin routes under `/api/v1`:

- `POST /films` - create a film`
- `PUT /films/:id` - update a film
- `DELETE /films/:id` - delete a film

The protected routes require authentication and the admin-only middleware.

## Project Structure

- `backend/main.go` - application entrypoint
- `backend/config/database.go` - database connection and migration setup
- `backend/internal/domain` - domain models and interfaces
- `backend/internal/handler` - HTTP handlers
- `backend/internal/service` - business logic
- `backend/internal/repository` - data access layer
- `backend/pkg/middleware` - auth and rate-limiting middleware
- `backend/pkg/response` - shared API response helpers
- `backend/routes/routes.go` - route registration

## Notes

- The backend enables CORS for `http://localhost:5174`.
- Database tables for users and films are migrated automatically when the app starts.
- If you want to test the API quickly, start with `POST /register`, then `POST /login`, and use the returned token for protected endpoints.

## Running Build

To produce a binary:

```bash
go build -o server.exe .
```

Then run it with:

```bash
server.exe
```
