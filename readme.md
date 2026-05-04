**Web Streaming (Backend)**

A small, beginner-friendly backend for a web streaming example. This repo is intended for learning how a Go web service is organized and run.

✅ Clean architecture
✅ JWT auth middleware
✅ Admin only middleware
✅ Input validation
✅ Duplicate check (email & username)
✅ Rate limiting
✅ CORS
✅ Zerolog logging
✅ Consistent response format
✅ /api/v1 versioning
✅ json:"-" pada password
✅ GORM tags pada Filem

**Prerequisites**
- Install Go (1.18+ recommended).

**Quick Start**
1. Open a terminal and go to the backend folder:

	```bash
	cd backend
	```

2. Download dependencies and run the server:

	```bash
	go mod tidy
	go run .
	```

3. The server will print the listening port in the console — open that address in your browser or use `curl`/Postman to call the routes.

**Project Structure**
- **main:** [backend/main.go](backend/main.go) — application entrypoint and server start.
- **config:** [backend/config/database.go](backend/config/database.go) — configuration helpers (database setup).
- **internal/domain:** [backend/internal/domain](backend/internal/domain) — domain models (`film.go`, `user.go`).
- **internal/handler:** [backend/internal/handler](backend/internal/handler) — HTTP handlers for routes.
- **internal/service:** [backend/internal/service](backend/internal/service) — business logic.
- **internal/repository:** [backend/internal/repository](backend/internal/repository) — data access layer.
- **pkg/middleware:** [backend/pkg/middleware](backend/pkg/middleware) — middleware utilities (auth, etc.).
- **routes:** [backend/routes/routes.go](backend/routes/routes.go) — route definitions and wiring.

**How to Explore (for Beginners)**
- To see how a request flows: follow the handler in `internal/handler`, then the service in `internal/service`, and finally repository in `internal/repository`.
- Look at `routes/routes.go` to see which endpoints are available and which handler each one calls.

**Next Steps / Tips**
- Add environment variables or a `.env` file if the project expects database credentials.
- Run `go build` to produce an executable: `go build -o server .` then run `./server`.

If you want, I can add example curl commands or a tiny quickstart that shows one working endpoint.

