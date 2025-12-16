# Sistem Pelaporan Prestasi Mahasiswa (Backend API)

**Proyek Akhir – Pemrograman Backend Lanjut (Praktikum)**
**DIV Teknik Informatika – Universitas Airlangga**

---

## 👤 Tentang Saya

| Atribut | Detail                |
| ------- | --------------------- |
| Nama    | **Agnesti Wulansari** |
| NIM     | **434231002**         |
| Kelas   | **TI-C1**             |

---

## Tentang Proyek

Sistem Pelaporan Prestasi Mahasiswa adalah layanan **Backend berbasis REST API** yang dirancang untuk memfasilitasi **pelaporan, pengelolaan, dan verifikasi prestasi mahasiswa**.

Sistem ini menggunakan **arsitektur Hybrid Database (Polyglot Persistence)** untuk menangani:

* **Data relasional yang terstruktur** (User, Role, Permission)
* **Data prestasi yang dinamis** (detail lomba, dokumen, dan atribut fleksibel)

Selain itu, sistem ini mengimplementasikan **Role-Based Access Control (RBAC)** untuk memastikan keamanan dan pembagian hak akses antara **Mahasiswa, Dosen Wali, dan Admin**.

---

## Fitur Utama

### 🔐 1. Autentikasi & Otorisasi (RBAC)

* Login aman menggunakan **JWT (JSON Web Token)**
* Middleware otorisasi berbasis **Permission**
  *(contoh: `achievement:create`, `achievement:verify`)*
* Manajemen user berdasarkan role:

  * Admin
  * Dosen Wali
  * Mahasiswa

---

### 🏆 2. Manajemen Prestasi (Hybrid Storage)

* Mahasiswa dapat **menginput prestasi** dengan field dinamis:

  * Akademik
  * Kompetisi
  * Organisasi
  * dll
* **Upload file bukti prestasi** (Sertifikat / SK)
* **Workflow Status Prestasi**:

  ```text
  Draft → Submitted → Verified / Rejected
  ```

---

### ✅ 3. Verifikasi Prestasi (Dosen Wali)

* Melihat daftar prestasi mahasiswa bimbingan
* Melakukan:

  * Approve (Verifikasi)
  * Reject (Penolakan) dengan catatan revisi

---

### 📊 4. Pelaporan & Analitik

* Statistik total prestasi per periode
* **Leaderboard** mahasiswa berprestasi

---

## 🛠️ Teknologi yang Digunakan

| Kategori     | Teknologi       | Kegunaan                                |
| ------------ | --------------- | --------------------------------------- |
| Language     | Go (Golang)     | Bahasa pemrograman utama                |
| Framework    | Gin Gonic       | HTTP Web Framework cepat                |
| RDBMS        | PostgreSQL      | User, Role, Permission, Relasi Akademik |
| NoSQL        | MongoDB         | Detail prestasi fleksibel (JSON)        |
| ORM          | GORM            | ORM PostgreSQL                          |
| Mongo Driver | mongo-go-driver | Driver resmi MongoDB                    |
| Dokumentasi  | Swagger         | Dokumentasi API otomatis                |
| Auth         | JWT             | Token-based Authentication              |

Referensi resmi:

* [https://go.dev/](https://go.dev/)
* [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [https://www.postgresql.org/](https://www.postgresql.org/)
* [https://www.mongodb.com/](https://www.mongodb.com/)
* [https://gorm.io/](https://gorm.io/)
* [https://github.com/mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
* [https://github.com/swaggo/swag](https://github.com/swaggo/swag)

---

## 🏗️ Arsitektur Database

Sistem ini menggunakan pendekatan **Polyglot Persistence**:

### 1️⃣ PostgreSQL (Data Struktural)

Digunakan untuk data dengan relasi ketat:

* `users`
* `roles`
* `permissions`
* `students`
* `lecturers`
* `achievement_references`
  *(menyimpan status & relasi ke MongoDB)*

### 2️⃣ MongoDB (Data Fleksibel)

Digunakan untuk data prestasi yang strukturnya dinamis:

* Collection `achievements`
* Struktur dokumen dapat berbeda tergantung jenis prestasi
  *(Juara Lomba, Publikasi Jurnal, Organisasi, dll)*

---

## ⚙️ Cara Instalasi & Menjalankan

### 🔹 Prasyarat

* Go **v1.20+**
* PostgreSQL Server
* MongoDB Server

---

### 🔹 Langkah-langkah

#### 1️⃣ Clone Repository

```bash
git clone https://github.com/agnestiw/sistem-pelaporan-prestasi-mahasiswa.git
cd sistem-pelaporan-prestasi-mahasiswa
```

#### 2️⃣ Konfigurasi Environment

Buat file `.env` di root folder:

```env
# Server Config
PORT=8080
APP_ENV=development

# PostgreSQL Config
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password_anda
DB_NAME=prestasi_db
DB_PORT=5432
DB_SSLMODE=disable

# MongoDB Config
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=prestasi_mongo_db

# JWT Config
JWT_SECRET=rahasia_super_aman_123
JWT_TTL=24
```

#### 3️⃣ Install Dependencies

```bash
go mod tidy
```

#### 4️⃣ Jalankan Aplikasi

```bash
go run main.go
```

#### 5️⃣ Akses Dokumentasi API

```text
http://localhost:8080/swagger/index.html
```

---

## 📡 Daftar Endpoint Utama

| Method | Endpoint                          | Deskripsi           | Akses         |
| ------ | --------------------------------- | ------------------- | ------------- |
| POST   | `/api/v1/auth/login`              | Login & ambil token | Public        |
| POST   | `/api/v1/achievements`            | Input prestasi      | Mahasiswa     |
| POST   | `/api/v1/achievements/:id/submit` | Submit prestasi     | Mahasiswa     |
| POST   | `/api/v1/achievements/:id/verify` | Verifikasi prestasi | Dosen Wali    |
| POST   | `/api/v1/achievements/:id/reject` | Tolak prestasi      | Dosen Wali    |
| GET    | `/api/v1/reports/statistics`      | Statistik prestasi  | Admin / Dosen |

---

## 🧪 Pengujian (Testing)

Menjalankan unit test:

```bash
go test ./tests/... -v
```

---