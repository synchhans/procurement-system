# Procurement System

Sistem Pengadaan Barang berbasis website (Technical Test).

## Tech Stack

### Backend
- **Language:** Go (Golang) 1.18+
- **Framework:** [Go Fiber](https://gofiber.io/) (Framework web yang cepat & ringan)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Auth:** JWT (JSON Web Token) & Bcrypt
- **Architecture:** Handlers, Models, Middleware

### Frontend
- **Library:** jQuery (DOM & AJAX)
- **Styling:** Tailwind CSS (via CDN)
- **Components:** SweetAlert2 (Toast), FontAwesome

## Cara Menjalankan

### Prasyarat
- Go 1.18 atau lebih baru.
- PostgreSQL sudah terinstall dan berjalan.

### 1. Clone Repository
```bash
git clone https://github.com/username-anda/procurement-system.git
cd procurement-system
```

### 2. Setup Database

Buat database kosong di PostgreSQL bernama `procurement_db`.

### 3. Konfigurasi Environment

Buat file `.env` di root folder project, salin kode di bawah ini:

```bash
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=procurement_db
DB_PORT=5432
JWT_SECRET=rahasia_super_aman
PORT=3000
WEBHOOK_URL=https://webhook.site/uuid
```

### 4. Isi Data Dummy (Seeding)

```bash
go mod tidy
go run cmd/seed/main.go
```

### 5. Jalankan Backend

```bash
go run cmd/api/main.go
```

Server akan berjalan di http://localhost:3000

### 6. Jalankan Frontend

Dapat membuka file `public/index.html` secara langsung atau menggunakan *Live Server*.

### 7. Akun Demo

Gunakan akun ini untuk login:

username: `admin`<br>
password: `password123`