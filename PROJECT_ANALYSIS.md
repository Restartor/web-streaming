# 1. Project Overview
- **Tujuan project**: Backend API platform web streaming film dengan fitur auth, katalog film, watchlist, dan riwayat tontonan.
- **Teknologi yang digunakan**:
  - Go 1.25
  - Gin (HTTP framework)
  - GORM + PostgreSQL
  - JWT (access token)
  - Refresh token (disimpan di DB)
  - Zerolog
  - Ulule limiter (rate limit)
- **Architecture pattern**: Layered architecture `handler -> service -> repository -> database`, dengan kontrak interface di layer `internal/domain`.
- **Current development state**:
  - Fitur inti backend sudah berjalan untuk auth dasar + film + watchlist + history.
  - Belum ada testing otomatis.
  - Belum ada fitur production-critical seperti email verification, forgot password, observability lengkap, subscription/payment, recommendation.

## Flow aplikasi (high-level)
1. Request masuk ke Gin route group `/api/v1`.
2. Middleware dijalankan (rate limiter, auth JWT, admin check sesuai route).
3. Handler memvalidasi input dan memanggil service.
4. Service menjalankan business logic.
5. Repository akses PostgreSQL via GORM.
6. Response distandardisasi ke format `{ data, error }`.

# 2. Folder Structure
## Root
- `backend/`: source code backend utama.
- `PROJECT_ANALYSIS.md`: dokumen analisis ini.

## backend/internal/
Berisi core business layer:
- `internal/domain/`: entity/model + interface contract repository/service.
- `internal/handler/`: HTTP handler (binding request, response API).
- `internal/service/`: business logic per modul.
- `internal/repository/`: implementasi data access ke DB.

## backend/internal/handler/
- `user_handler.go`: register, login, refresh-token, logout.
- `film_handler.go`: list/search film + admin CRUD film.
- `watchlist_handler.go`: add/remove/get watchlist user.
- `watched_handler.go`: get/delete/record history user.

## backend/internal/service/
- `user_service.go`: register/login, generate JWT, refresh token rotation, logout.
- `film_service.go`: operasi film.
- `watchlist_service.go`: operasi watchlist.
- `watched_service.go`: operasi history watch.

## backend/internal/repository/
- `user_repository.go`: query user.
- `refresh_token_repository.go`: simpan/cari/hapus refresh token.
- `film_repository.go`: query film + pagination + title search.
- `watchlist_repository.go`: query watchlist.
- `watched_repository.go`: query history dan upsert last watched.

## backend/pkg/
Shared utilities:
- `pkg/middleware/auth.middleware.go`: validasi Bearer JWT, inject `user_id`, `role`, `username` ke context.
- `pkg/middleware/rate_limiter.go`: rate limit per endpoint.
- `pkg/adminOnly.go`: pembatasan role admin.
- `pkg/response/response.go`: response wrapper JSON.
- `pkg/logger/logger.go`: inisialisasi logger.

## backend/routes/
- `routes.go`: definisi semua route public/auth/admin + middleware chain.

## backend/config/
- `database.go`: koneksi DB, automigrate tabel domain.

## domain/entity/model
Entity utama berada di:
- `internal/domain/user.go`
- `internal/domain/film.go`

# 3. Authentication System
## Login flow
1. `POST /api/v1/login` menerima `email`, `password`.
2. Service cari user by email.
3. Verifikasi password dengan bcrypt.
4. Generate access token JWT (`JWT_SECRET`, durasi `ACCESS_TOKEN_DURATION`).
5. Generate refresh token UUID, simpan ke tabel `refresh_tokens` dengan expiry (`REFRESH_TOKEN_DURATION`).
6. Return `access_token` + `refresh_token`.

## Register flow
1. `POST /api/v1/register` menerima username/email/password.
2. Service set default role = `user`.
3. Cek duplikasi email dan username.
4. Hash password dengan bcrypt.
5. Simpan user ke DB.

## JWT usage
- JWT dipakai di protected routes.
- Claims yang disimpan: `user_id`, `username`, `role`, `exp`.
- Auth middleware membaca header `Authorization: Bearer <token>` dan parse token.

## Refresh token
- Refresh token bersifat stateful (disimpan di DB).
- Endpoint `POST /api/v1/refresh-token`:
  - Validasi token ada di DB dan belum expired.
  - Generate access token baru.
  - Rotate refresh token (hapus lama, simpan baru).

## Middleware auth
- `AuthMiddleware()`:
  - validasi header Authorization.
  - parse JWT.
  - set context: `user_id`, `role`, `username`.

## Role system
- Role disimpan pada `User.Role`.
- Role default saat register: `user`.
- Endpoint admin memakai `AdminOnly()` (cek context `role == "admin"`).

## Email verification
- **Belum ada** implementasi.

## Forgot password
- **Belum ada** implementasi.

# 4. Entity / Domain Analysis
## User
Fields:
- `id` (uint)
- `username` (string)
- `email` (string)
- `password` (string, tidak di-serialize ke JSON)
- `role` (string)
- `created_at` (time)

Digunakan oleh:
- auth register/login/logout
- role authorization admin

## RefreshToken
Fields:
- `id` (uint, primary key)
- `user_id` (uint, indexed)
- `token` (string, unique)
- `expires_at` (time)
- `created_at` (time)

Digunakan oleh:
- login (create)
- refresh token rotation
- logout (delete by user)

## Filem (Film)
Fields:
- `id` (uint)
- `title` (string)
- `description` (string)
- `genre` (`pq.StringArray` / text[])
- `year` (int)
- `poster_url` (string)
- `rating` (float64)
- `video_url` (string)

Digunakan oleh:
- listing/search film public
- CRUD film admin
- referensi watchlist/history

## UserWatchList
Fields:
- `user_id` (uint, PK)
- `film_id` (uint, PK)

Digunakan oleh:
- add/remove/get watchlist user login

## UserHistory
Fields:
- `user_id` (uint, PK)
- `film_id` (uint, PK)
- `last_watched_at` (time)

Digunakan oleh:
- record watch (upsert)
- get/delete history

## Supporting domain types
- `RegisterInput`, `LoginInput`
- `PaginationQuery`, `PaginatedFilms`
- Interface contracts `UserRepository`, `UserService`, `FilmRepository`, `FilmService`, dll.

# 5. API Endpoint Analysis
## Auth
| Method | Endpoint | Access | Description | Handler |
|---|---|---|---|---|
| POST | /api/v1/register | Public | Registrasi user baru | UserHandler.Register |
| POST | /api/v1/login | Public | Login dan issue access+refresh token | UserHandler.Login |
| POST | /api/v1/refresh-token | Public | Refresh access token + rotate refresh token | UserHandler.RefreshToken |
| POST | /api/v1/logout | Authenticated | Logout, revoke semua refresh token user | UserHandler.Logout |

## Movie / Film
| Method | Endpoint | Access | Description | Handler |
|---|---|---|---|---|
| GET | /api/v1/films | Public | Ambil daftar film (pagination) | FilmHandler.GetAllFilms |
| GET | /api/v1/films/search?title=... | Public | Cari film berdasarkan title | FilmHandler.GetFilmByTitle |
| POST | /api/v1/films | Admin | Tambah film | FilmHandler.CreateFilm |
| PUT | /api/v1/films/:id | Admin | Update film | FilmHandler.UpdateFilm |
| DELETE | /api/v1/films/:id | Admin | Hapus film | FilmHandler.DeleteFilm |

## Watchlist
| Method | Endpoint | Access | Description | Handler |
|---|---|---|---|---|
| GET | /api/v1/watchlist | Authenticated | Ambil watchlist user | WatchlistHandler.GetWatchlist |
| POST | /api/v1/watchlist | Authenticated | Tambah film ke watchlist | WatchlistHandler.AddToWatchlist |
| DELETE | /api/v1/watchlist/:id | Authenticated | Hapus film dari watchlist | WatchlistHandler.RemoveFromWatchlist |

## History
| Method | Endpoint | Access | Description | Handler |
|---|---|---|---|---|
| GET | /api/v1/history | Authenticated | Ambil history user | HistoryHandler.GetAllHistory |
| POST | /api/v1/history/:id | Authenticated | Record watch event untuk film id | HistoryHandler.RecordWatch |
| DELETE | /api/v1/history/:id | Authenticated | Hapus satu history film | HistoryHandler.DeleteHistoryOne |
| DELETE | /api/v1/history | Authenticated | Hapus semua history user | HistoryHandler.DeleteAllHistory |

# 6. Existing Features
- User register.
- User login.
- Access token JWT.
- Refresh token persistence + rotation.
- Logout (revoke refresh token by user).
- Public film listing dengan pagination.
- Film search by title.
- Admin-only film CRUD.
- Watchlist add/remove/list per user.
- Watch history record/list/delete.
- Role-based route protection (admin/user).
- Rate limiting per endpoint.
- Standardized API response wrapper.
- Auto migration schema saat startup.

# 7. Missing Features
Kebutuhan yang belum ada untuk production-ready streaming platform:
- Email verification flow.
- Forgot/reset password flow.
- Secure refresh token hardening (hash token di DB, device/session metadata, selective logout per device).
- OAuth/social login.
- Health check endpoint + readiness/liveness.
- Observability lengkap (metrics, tracing, audit log).
- API documentation formal (OpenAPI/Swagger).
- Automated tests (unit/integration/e2e).
- Caching layer (Redis) untuk catalog/search.
- Recommendation system.
- Continue watching progress detail (position/duration).
- Subtitle/multi-audio support metadata.
- Content category/featured/trending endpoints.
- Ads management module (placement, targeting, frequency cap).
- Payment/subscription/billing module.
- Entitlement access control untuk premium content.
- File upload/asset management (poster, video processing).
- Advanced search/filter/sort endpoint.
- Pagination untuk watchlist/history (saat ini belum).
- Security hardening tambahan (CSRF strategy jika pakai cookie, brute-force protection account-level, secret management).

## Catatan TODO / unfinished dari source code
- `HistoryHandler.RecordWatch` melakukan bind body `film_id` tetapi nilai tersebut tidak dipakai (yang dipakai `:id` path param).
- Tidak ada endpoint/profile module walau context username sudah di-inject middleware.
- Tidak ada modul admin selain film CRUD.

# 8. Frontend Page Recommendation
| Page | Route | Purpose | Required API |
|---|---|---|---|
| Landing/Home | `/` | Menampilkan katalog film awal | `GET /api/v1/films`, `GET /api/v1/films/search` |
| Login | `/login` | Autentikasi user | `POST /api/v1/login` |
| Register | `/register` | Registrasi user baru | `POST /api/v1/register` |
| Film List | `/movies` | Browse film dengan pagination | `GET /api/v1/films` |
| Film Search Result | `/movies/search` | Tampilkan hasil pencarian title | `GET /api/v1/films/search` |
| Watchlist | `/watchlist` | Kelola daftar tontonan user | `GET /api/v1/watchlist`, `POST /api/v1/watchlist`, `DELETE /api/v1/watchlist/:id` |
| Watch History | `/history` | Riwayat tontonan user | `GET /api/v1/history`, `DELETE /api/v1/history/:id`, `DELETE /api/v1/history` |
| Player/Watch | `/watch/:id` | Playback + trigger watch record | `POST /api/v1/history/:id` |
| Admin Film Management | `/admin/films` | CRUD film untuk admin | `GET /api/v1/films`, `POST /api/v1/films`, `PUT /api/v1/films/:id`, `DELETE /api/v1/films/:id` |
| Session Refresh (silent) | internal guard | Menjaga sesi login tetap aktif | `POST /api/v1/refresh-token`, `POST /api/v1/logout` |

# 9. UI/UX Recommendation
- **Halaman prioritas Figma**:
  1. Login/Register
  2. Home/Film List
  3. Movie Detail + Player
  4. Watchlist
  5. History
  6. Admin Film Dashboard
- **Reusable component**:
  - Navbar (guest/user/admin state)
  - MovieCard, MovieGrid, Pagination
  - SearchBar
  - AuthForm (shared login/register fields)
  - ProtectedRoute + RoleGuard
  - Toast/Alert API error
  - Confirmation modal (delete watchlist/history/film)
- **Layout recommendation**:
  - Public: hero + rail/grid katalog.
  - Authenticated: personalized rails (watchlist/history shortcut).
  - Admin: table/grid CRUD + form drawer/modal.
- **Mobile-first consideration**:
  - Bottom navigation untuk page utama.
  - Card vertikal dan horizontal scroll section.
  - Form auth satu kolom.
  - Aksi penting mudah dijangkau (thumb zone).
- **Ads placement strategy (future module)**:
  - Home feed: in-feed card slot berkala.
  - Player: pre-roll/mid-roll marker (butuh backend ads service).
  - Hindari gangguan di auth flow.
- **Streaming platform UX pattern**:
  - Continue Watching rail.
  - My List/Watchlist persist across sessions.
  - Search cepat dengan debounce.
  - Fast resume dari history.

# 10. Backend-Frontend Contract Draft
## Home Page
- `GET /api/v1/films?page=1&limit=10`
- `GET /api/v1/films/search?title=<query>` (untuk search global)

## Login Page
- `POST /api/v1/login`
- Jika access token expired: `POST /api/v1/refresh-token`

## Register Page
- `POST /api/v1/register`

## Film Detail / Player Page
- (Saat ini belum ada endpoint detail by ID)
- Playback tracking:
  - `POST /api/v1/history/:id`
- Action user:
  - `POST /api/v1/watchlist`
  - `DELETE /api/v1/watchlist/:id`

## Watchlist Page
- `GET /api/v1/watchlist`
- `POST /api/v1/watchlist`
- `DELETE /api/v1/watchlist/:id`

## History Page
- `GET /api/v1/history`
- `DELETE /api/v1/history/:id`
- `DELETE /api/v1/history`

## Admin Film Management Page
- `GET /api/v1/films`
- `POST /api/v1/films`
- `PUT /api/v1/films/:id`
- `DELETE /api/v1/films/:id`

## Session/Auth Guard (global frontend)
- Simpan `access_token` + `refresh_token` dari login.
- Saat 401 karena token expired, trigger `POST /api/v1/refresh-token`.
- Saat logout, panggil `POST /api/v1/logout` lalu clear local auth state.
