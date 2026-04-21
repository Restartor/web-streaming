# web-streaming

> Proyek ini dibuat untuk tujuan belajar (*study purposes only*).

Web Streaming adalah REST API sederhana yang dibangun menggunakan **Go** dengan pendekatan **Clean Architecture**. Proyek ini mendemonstrasikan bagaimana memisahkan logika bisnis, akses data, dan lapisan HTTP secara terstruktur agar mudah dipahami, diuji, dan dikembangkan.

---

## Daftar Isi

- [Struktur Folder](#struktur-folder)
- [Lapisan Arsitektur](#lapisan-arsitektur)
- [Endpoint API](#endpoint-api)
- [Environment Variables](#environment-variables)
- [Cara Menjalankan](#cara-menjalankan)
- [Cara Menjalankan Test](#cara-menjalankan-test)
- [Dependensi](#dependensi)

---

## Struktur Folder

```
/web-streaming
├── main.go                     ← titik masuk aplikasi
├── go.mod / go.sum             ← manajemen modul Go
├── config/
│   ├── database.go             ← konfigurasi koneksi database
│   └── redis.go                ← konfigurasi koneksi Redis
├── internal/
│   ├── domain/                 ← struct & interface (tidak bergantung ke lapisan lain)
│   │   ├── film.go
│   │   └── user.go
│   ├── repository/             ← implementasi akses data
│   │   ├── film_repository.go
│   │   └── user_repository.go
│   ├── service/                ← logika bisnis
│   │   ├── film_service.go
│   │   └── auth_service.go
│   └── handler/                ← menerima request HTTP, mengirim response
│       ├── film_handler.go
│       └── auth_handler.go
├── pkg/
│   ├── middleware/
│   │   └── auth.go             ← middleware autentikasi Bearer token
│   └── utils/
│       └── response.go         ← helper untuk menulis JSON response
└── routes/
    └── routes.go               ← mendaftarkan semua route HTTP
```

---

## Lapisan Arsitektur

Proyek ini mengikuti prinsip **Clean Architecture** dengan empat lapisan utama:

| Lapisan | Paket | Tanggung Jawab |
|---|---|---|
| **Domain** | `internal/domain` | Mendefinisikan struct entitas (`Film`, `User`) dan interface repository. Tidak bergantung ke lapisan lain. |
| **Repository** | `internal/repository` | Mengimplementasikan interface repository dari domain. Bertanggung jawab atas penyimpanan dan pengambilan data. |
| **Service** | `internal/service` | Mengandung logika bisnis. Hanya bergantung ke interface domain, bukan implementasi konkret. |
| **Handler** | `internal/handler` | Menerima request HTTP, memanggil service, dan mengembalikan response JSON. |

Arah dependensi: `Handler → Service → Repository (interface) ← Repository (implementasi)`

---

## Endpoint API

### `GET /films`

Mengambil daftar semua film.

**Response `200 OK`:**
```json
[
  { "id": 1, "title": "Inception" }
]
```

---

### `POST /login`

Melakukan autentikasi pengguna.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "rahasia"
}
```

**Response `200 OK`:**
```json
{ "success": true }
```

**Response `401 Unauthorized`:**
```json
{ "error": "invalid credentials" }
```

---

## Environment Variables

| Variabel | Default | Keterangan |
|---|---|---|
| `PORT` | `8080` | Port tempat server HTTP berjalan |
| `WEB_STREAMING_AUTH_TOKEN` | *(wajib diisi)* | Token Bearer yang divalidasi oleh middleware `RequireAuth` |

---

## Cara Menjalankan

**Prasyarat:** Go 1.22 atau lebih baru.

```bash
# 1. Clone repository
git clone https://github.com/Restartor/web-streaming.git
cd web-streaming

# 2. Install dependensi
go mod tidy

# 3. Jalankan server (port default: 8080)
go run .

# Atau tentukan port sendiri
PORT=9090 go run .
```

Server akan berjalan di `http://localhost:8080`.

---

## Cara Menjalankan Test

```bash
# Jalankan semua test
go test ./...

# Jalankan test dengan output verbose
go test -v ./...
```

---

## Dependensi

| Paket | Kegunaan |
|---|---|
| [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto) | Hashing password menggunakan **bcrypt** |
