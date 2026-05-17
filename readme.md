# Web Streaming Backend — Guide untuk Frontend, Testing, dan DevOps

Panduan lengkap (ringkas dan praktis) untuk `web-streaming` backend agar tim frontend, tim testing, dan DevOps/infra bisa langsung pakai dan deploy ke produksi.

Catatan: struktur kode mengikuti pola handler → service → repository, gunakan file di `backend/` untuk menjalankan server.

## Ringkasan Singkat
- Bahasa & framework: Go + Gin
- ORM: GORM (Postgres driver)
- Auth: JWT (access + refresh tokens)
- Logging: Zerolog (JSON stdout)
- Rate limiting: ulule/limiter
- Password hashing: bcrypt

Tujuan: sediakan API REST untuk registrasi/login, katalog film publik, watchlist/history pengguna, dan manajemen film (admin).

---

## Untuk Frontend — Integrasi Praktis

- Base URL: `{{BASE_URL}}/api/v1` (konfigurasi CORS lewat `ALLOWED_ORIGINS` di env).
- Auth flow:
  - `POST /register` — body: `{ username, email, password }` → response wrapper `{ data, error }`.
  - `POST /login` — body: `{ email, password }` → `data` berisi `{ access_token, refresh_token, expires_in }`.
  - `POST /refresh-token` — body: `{ refresh_token }` → returns new `access_token` (possible new refresh token).
  - `POST /logout` — authenticated: revokes refresh tokens.

- Sending tokens:
  - Use header: `Authorization: Bearer <access_token>` for protected endpoints.
  - Recommended storage for refresh token (browser): HttpOnly Secure cookie. Access token in memory (client-side) and refreshed frequently.

- Endpoints penting:
  - Public: `GET /films?page=&limit=`, `GET /films/search?title=`.
  - Authenticated: `GET/POST/DELETE /watchlist`, `GET/DELETE /history`.
  - Admin-only: `POST/PUT/DELETE /films` (requires `role=admin` claim in JWT).

- Error handling: Semua response mengikuti format `{ data, error }`. Jika `error != null`, tampilkan pesan kepada user dari `error.message`.

- Pagination: query params `page` (1-based) dan `limit`. UI harus menampilkan total pages/total items bila available.

- Caching & UX:
  - Cache film list (short TTL 1–5 min). Invalidate cache on admin create/update/delete.
  - Debounce search input and cache recent queries.

Contoh: Login (curl)
```bash
curl -X POST {{BASE_URL}}/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"pass123"}'
```

Contoh response (sederhana):
```json
{
  "data": { "access_token": "...", "refresh_token": "...", "expires_in": 3600 },
  "error": null
}
```

---

## Untuk Tim Testing

- Unit tests:
  - Fokus di `internal/service` (mock repository interfaces). Gunakan table-driven tests.
  - Gunakan mocking (manual mock atau mock generator) agar tidak bergantung DB.

- Integration tests:
  - Jalankan terhadap Postgres ephemeral (Docker). Seed minimal data (admin, user, beberapa film).
  - Gunakan `go test ./...` dengan build tag `integration` untuk memisahkan.

- E2E / API tests:
  - Gunakan Postman/Newman, atau test runner yang memanggil endpoints pada staging.
  - Test flows: register → login → add watchlist → get watchlist → logout.

- Security tests:
  - Verifikasi role checks (admin-only endpoints), JWT expiry behavior, password hashing.
  - Pastikan refresh token reuse is handled (revoked after logout if implemented).

- Load testing:
  - Gunakan `k6` atau `JMeter` untuk endpoints read-heavy (`GET /films`). Simulasikan sustained users and burst traffic.

- CI tips:
  - Run `gofmt`, `golangci-lint`, `go vet`, and `go test ./...` in CI.
  - For integration tests, spin up a Postgres container in CI job and run migrations.

---

## Untuk DevOps / Produksi (Langkah demi langkah)

1. Secrets & Config
   - Simpan `JWT_SECRET`, `REFRESH_TOKEN_SECRET`, DB credentials di secret manager (Vault, AWS Secrets Manager, Azure Key Vault, Kubernetes Secret).
   - Jangan cek-in `.env` ke repo.

2. Database migrations
   - Jangan andalkan `AutoMigrate` untuk produksi. Gunakan migration tool (e.g., `golang-migrate`).
   - CI/CD harus menerapkan migration job pra-deploy atau migration job terkontrol di fase release.

3. Containerization
   - Buat `Dockerfile` multi-stage: build binary, lalu copy ke minimal runtime image (scratch atau alpine).
   - Tambah `docker-compose.yml` untuk local dev (app + postgres + adminer).

4. Health & Probes
   - Tambahkan `/health` endpoint (liveness) dan readiness (DB connection ok).
   - Konfigurasikan probes di K8s manifests.

5. TLS & Ingress
   - Terminate TLS di ingress (NGINX, ALB) atau Cloud Load Balancer.

6. Observability
   - Output Zerolog JSON ke stdout.
   - Kirim metrics (Prometheus) dan traces (OpenTelemetry) jika memungkinkan.
   - Centralized logs (ELK/Elastic/CloudWatch).

7. Rate limiting & scaling
   - For distributed rate limits, use Redis-backed limiter or global gateway.
   - App is stateless: scale pods horizontally behind a LB.

8. Backups & DR
   - Schedule DB backups, test restore procedures.

9. CI/CD pipeline (recommended steps)
   - Lint & static analysis
   - Unit tests
   - Build Docker image
   - Run integration tests against ephemeral infra
   - Publish artifact to registry
   - Deploy to staging, run smoke tests
   - Deploy to production with canary/blue-green

---

## Snippets yang berguna

Dockerfile (contoh minimal):

```dockerfile
# build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/web-streaming ./backend

# runtime stage
FROM alpine:3.18
COPY --from=builder /bin/web-streaming /bin/web-streaming
ENV GIN_MODE=release
EXPOSE 1010
ENTRYPOINT ["/bin/web-streaming"]
```

docker-compose (dev, ringkas):

```yaml
version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: web_streaming
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - 5432:5432
  app:
    build: .
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=web_streaming
    ports:
      - 1010:1010
    depends_on:
      - db
```

CI snippet (GitHub Actions) — ringkas:

```yaml
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - name: Install deps
        run: go mod download
      - name: Lint
        run: golangci-lint run ./...
      - name: Test
        run: go test ./...
      - name: Build
        run: go build -o web-streaming ./backend
```

---

## Environment variables (ringkasan)

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `ALLOWED_ORIGINS` (CORS)
- `JWT_SECRET`, `REFRESH_TOKEN_SECRET`
- `ACCESS_TOKEN_DURATION`, `REFRESH_TOKEN_DURATION`

Pastikan infrastruktur menyuntikkan value ini dari secret manager.

---

## Checklist produksi (ceklis terperinci untuk siap deploy)

- [ ] Secrets terenkripsi di secret manager
- [ ] Migrations disiapkan dan diuji pada staging
- [ ] Health/readiness probes terpasang
- [ ] TLS dan CORS dikonfigurasi untuk origin frontend
- [ ] Centralized logging & metrics tersedia
- [ ] Backups DB diatur
- [ ] Autoscaling & rate limiting dikonfigurasi
- [ ] CI tests (unit + integration) lulus

---

## Kontak & kontribusi

Untuk perubahan API, update dokumentasi OpenAPI / Postman collection dan informasikan tim frontend.

Jika mau, saya bisa:
- Buat `Dockerfile` / `docker-compose.yml` (sudah contoh di atas)
- Tambahkan OpenAPI / Swagger spec
- Buat GitHub Actions CI lengkap dengan integration test using Postgres container

Ketik mana yang mau saya kerjakan selanjutnya.
