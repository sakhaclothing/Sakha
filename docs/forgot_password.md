# Forgot Password API Documentation

## Overview

Sistem forgot password memungkinkan user untuk mereset password mereka melalui email. Sistem ini menggunakan token yang aman dan memiliki expiry time.

## Endpoints

### 1. Forgot Password Request

**POST** `/auth/forgot-password`

Mengirim request untuk reset password berdasarkan email.

#### Request Body

```json
{
  "email": "user@example.com"
}
```

#### Response

```json
{
  "message": "Link reset password telah dikirim ke email Anda",
  "reset_link": "https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/reset-password?token=abc123..."
}
```

#### Error Responses

- `400` - Email tidak boleh kosong
- `500` - Gagal membuat token reset / Gagal menyimpan token reset

### 2. Reset Password

**POST** `/auth/reset-password`

Menggunakan token untuk mereset password.

#### Request Body

```json
{
  "token": "abc123...",
  "password": "newpassword123"
}
```

#### Response

```json
{
  "message": "Password berhasil direset"
}
```

#### Error Responses

- `400` - Token tidak boleh kosong / Password tidak boleh kosong / Password minimal 6 karakter / Token reset tidak valid atau sudah digunakan / Token reset sudah expired
- `404` - User tidak ditemukan
- `500` - Gagal memeriksa token / Gagal hash password / Gagal update password

## Flow Diagram

```
1. User request forgot password dengan email
   ↓
2. Sistem validasi email ada di database
   ↓
3. Generate reset token (expired 1 jam)
   ↓
4. Invalidate token lama untuk email tersebut
   ↓
5. Simpan token baru ke database
   ↓
6. Kirim email dengan link reset
   ↓
7. User klik link dan input password baru
   ↓
8. Sistem validasi token (valid, belum expired, belum used)
   ↓
9. Update password user
   ↓
10. Mark token sebagai used
```

## Security Features

1. **Token Expiry**: Token reset expired dalam 1 jam
2. **Single Use**: Token hanya bisa digunakan sekali
3. **Email Validation**: Hanya email yang terdaftar yang bisa request reset
4. **Token Invalidation**: Token lama di-invalidate saat request baru
5. **Secure Token Generation**: Menggunakan crypto/rand untuk generate token

## Email Configuration

Untuk production, set environment variables berikut:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

## Database Schema

### Collection: `password_resets`

```json
{
  "_id": "ObjectId",
  "email": "user@example.com",
  "token": "abc123...",
  "expires_at": "2024-01-01T12:00:00Z",
  "used": false,
  "created_at": "2024-01-01T11:00:00Z"
}
```

## Testing

### Test Case 1: Request Reset Password

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

### Test Case 2: Reset Password

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token": "abc123...", "password": "newpassword123"}'
```

## Notes

1. Dalam development, email akan di-log ke console (mock email)
2. Reset link dalam response hanya untuk development, hapus di production
3. Password minimal 6 karakter
4. Token case-sensitive
5. Email validation case-insensitive
