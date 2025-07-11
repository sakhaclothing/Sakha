# Sakha Backend (Go Fiber) - Dokumentasi

## ✨ Fitur Utama

- User Authentication: Register, login, JWT, reset password, verifikasi email (OTP)
- Profile Management: Edit profil (username, fullname, email), ganti password, verifikasi email baru dengan OTP
- Role Management: Admin dapat mengubah role user
- Tracker: Menyimpan dan menghitung visitor
- Terintegrasi dengan MongoDB

---

## 📁 Struktur Folder

- `controller/` — Semua handler endpoint (auth, tracker, dsb)
- `model/` — Struct data (User, EmailVerification, dsb)
- `route/` — Definisi semua route (endpoint)
- `config/` — Konfigurasi database, email, dsb
- `utils/` — Utility (hash password, token, email, OTP)
- `main.go` — Entry point (untuk Google Cloud Function)

---

## 🚀 Cara Menjalankan Lokal

1. **Clone repo:**
   ```bash
   git clone https://github.com/sakhaclothing/sakha-backend.git
   cd sakha-backend
   ```
2. **Set environment variable** (lihat contoh di `.env.example` jika ada)
3. **Install dependency:**
   ```bash
   go mod tidy
   ```
4. **Jalankan server:**
   ```bash
   go run main.go
   ```
   atau deploy ke Google Cloud Function sesuai instruksi di README utama.

---

## 🔑 Endpoint Penting

### Auth & Profile

- `POST /auth/register` — Register user baru (dengan OTP email)
- `POST /auth/login` — Login user (JWT)
- `POST /auth/profile` — Get profil user (dari token)
- `PUT /auth/profile` — Update profil (username, fullname, email)
  - Jika email diubah, sistem akan mengirim OTP ke email baru, dan user harus verifikasi.
- `POST /auth/verify-email` — Verifikasi email baru dengan OTP
- `PUT /user/password` — Ganti password (dengan validasi password lama)

### Tracker

- `POST /tracker` — Simpan data visitor
- `GET /tracker/count` — Ambil jumlah visitor

### Lainnya

- `POST /auth/forgot-password` — Request reset password (kirim email)
- `POST /auth/reset-password` — Reset password dengan token

---

## 📝 Alur Ganti Email & OTP

1. User update email via `PUT /auth/profile`
2. Backend:
   - Update email, set `is_verified: false`
   - Generate OTP, simpan ke koleksi `email_verifications`
   - Kirim OTP ke email baru
3. User submit OTP via `POST /auth/verify-email`
4. Backend:
   - Verifikasi OTP, set `is_verified: true` jika benar

---

## 🛡️ Keamanan

- Semua endpoint sensitif (profile, password, dsb) menggunakan JWT di header `Authorization`.
- Password di-hash dengan bcrypt.
- Email harus diverifikasi sebelum bisa digunakan penuh.

---

## 📄 Lisensi

Lihat file [LICENSE](../LICENSE).

---

**Catatan:**

- Untuk detail CI/CD dan deployment ke Google Cloud Function, lihat README utama di repo ini.
- Jika menambah fitur baru, update dokumentasi ini agar tim lain mudah memahami.
