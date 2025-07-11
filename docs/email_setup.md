# Email Setup Guide - Gmail SMTP

## Overview

Setup Gmail SMTP untuk mengirim email reset password yang sungguhan ke user.

## Step 1: Setup Gmail App Password

### 1.1 Enable 2-Factor Authentication

1. Buka [Google Account Settings](https://myaccount.google.com/)
2. Pilih "Security"
3. Aktifkan "2-Step Verification"

### 1.2 Generate App Password

1. Di halaman Security, cari "App passwords"
2. Klik "App passwords"
3. Pilih "Mail" dan "Other (Custom name)"
4. Masukkan nama: "Sakha Clothing API"
5. Klik "Generate"
6. **Copy password yang dihasilkan** (16 karakter)

## Step 2: Environment Variables

Set environment variables berikut di Google Cloud Functions:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=sakhaclothing@gmail.com
SMTP_PASSWORD=your_16_character_app_password
FROM_EMAIL=sakhaclothing@gmail.com
FROM_NAME=Sakha Clothing
```

### Cara Set di Google Cloud Functions:

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Pilih project Anda
3. Buka "Cloud Functions"
4. Pilih function "sakha"
5. Klik "Edit"
6. Scroll ke "Environment variables"
7. Tambahkan variables di atas

## Step 3: Test Email Sending

### Test dengan cURL:

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "zenkun.enterkill13@gmail.com"}'
```

### Expected Response:

```json
{
  "message": "Email sent"
}
```

## Step 4: Check Email

1. **Cek inbox** `zenkun.enterkill13@gmail.com`
2. **Cek spam folder** jika tidak ada di inbox
3. **Email akan berisi:**
   - Subject: "Reset Password - Sakha Clothing"
   - Link reset password yang bisa diklik
   - Token reset yang bisa di-copy

## Email Template

Email yang dikirim akan berisi:

```html
Subject: Reset Password - Sakha Clothing Halo, Anda telah meminta untuk mereset
password akun Sakha Clothing Anda. [RESET PASSWORD BUTTON] Atau copy paste link
berikut ke browser Anda:
https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/reset-password?token=abc123...
Penting: Link ini akan expired dalam 1 jam. Best regards, Sakha Clothing Team
```

## Troubleshooting

### 1. Email Tidak Terkirim

- ✅ Cek App Password sudah benar
- ✅ Cek 2FA sudah aktif
- ✅ Cek environment variables sudah set
- ✅ Cek logs di Cloud Functions

### 2. Email Masuk Spam

- ✅ Cek spam folder
- ✅ Mark email sebagai "Not Spam"
- ✅ Tambahkan sender ke contacts

### 3. SMTP Error

- ✅ Cek SMTP credentials
- ✅ Cek network connectivity
- ✅ Cek Gmail settings

### 4. App Password Error

- ✅ Regenerate App Password
- ✅ Pastikan 2FA aktif
- ✅ Gunakan password 16 karakter

## Security Notes

1. **App Password** lebih aman dari password biasa
2. **Environment variables** tidak akan terekspos di code
3. **Token reset** expired dalam 1 jam
4. **Email validation** memastikan email terdaftar

## Alternative Email Providers

Jika tidak ingin menggunakan Gmail, bisa gunakan:

### SendGrid:

```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your_sendgrid_api_key
```

### Mailgun:

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=your_mailgun_username
SMTP_PASSWORD=your_mailgun_password
```

## Testing Checklist

- [ ] 2FA Gmail aktif
- [ ] App Password dibuat
- [ ] Environment variables set
- [ ] Function deployed
- [ ] Test request sent
- [ ] Email received
- [ ] Link works
- [ ] Password reset successful
