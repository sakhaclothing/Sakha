# GitHub Secrets Setup Guide

## Overview

Setup GitHub Secrets untuk menyimpan credentials SMTP dan environment variables secara aman.

## Step 1: Setup GitHub Secrets

### 1.1 Buka Repository Settings

1. Buka repository GitHub Anda
2. Klik tab "Settings"
3. Scroll ke "Secrets and variables" → "Actions"
4. Klik "New repository secret"

### 1.2 Tambahkan Secrets Berikut

#### SMTP Configuration:

```
Name: SMTP_HOST
Value: smtp.gmail.com

Name: SMTP_PORT
Value: 587

Name: SMTP_USERNAME
Value: sakhaclothing@gmail.com

Name: SMTP_PASSWORD
Value: your_16_character_gmail_app_password

Name: FROM_EMAIL
Value: sakhaclothing@gmail.com

Name: FROM_NAME
Value: Sakha Clothing
```

#### Database & Security:

```
Name: MONGODB_URI
Value: your_mongodb_connection_string

Name: JWT_SECRET
Value: your_jwt_secret_key
```

#### Google Cloud Service Account:

```
Name: GCP_SA_KEY
Value: your_google_cloud_service_account_json_key
```

## Step 2: Setup Google Cloud Service Account

### 2.1 Create Service Account

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Pilih project Anda
3. Buka "IAM & Admin" → "Service Accounts"
4. Klik "Create Service Account"
5. Nama: "github-actions-deployer"
6. Description: "Service account for GitHub Actions deployment"

### 2.2 Assign Roles

Tambahkan roles berikut:

- `Cloud Functions Developer`
- `Service Account User`
- `Cloud Build Editor`

### 2.3 Create and Download Key

1. Klik service account yang baru dibuat
2. Tab "Keys" → "Add Key" → "Create new key"
3. Pilih "JSON"
4. Download file JSON
5. **Copy seluruh isi file JSON** ke GitHub Secret `GCP_SA_KEY`

## Step 3: GitHub Actions Workflow

Workflow sudah dibuat di `.github/workflows/deploy.yml` yang akan:

1. **Trigger otomatis** saat push ke main/master
2. **Setup Go environment**
3. **Authenticate ke Google Cloud**
4. **Deploy function dengan environment variables dari secrets**

## Step 4: Test Deployment

### 4.1 Push ke Main Branch

```bash
git add .
git commit -m "Add GitHub Actions workflow"
git push origin main
```

### 4.2 Monitor Deployment

1. Buka tab "Actions" di GitHub
2. Lihat workflow "Deploy to Google Cloud Functions"
3. Tunggu sampai selesai (biasanya 2-3 menit)

### 4.3 Test Function

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "zenkun.enterkill13@gmail.com"}'
```

## Security Benefits

### ✅ Keamanan GitHub Secrets:

- **Encrypted at rest** - Secrets dienkripsi di GitHub
- **Masked in logs** - Tidak akan muncul di logs
- **Access control** - Hanya admin yang bisa akses
- **Audit trail** - Log semua akses ke secrets

### ✅ Environment Variables:

- **Not in code** - Tidak ada credentials di source code
- **Per-environment** - Bisa beda untuk dev/staging/prod
- **Easy rotation** - Mudah ganti credentials

## Troubleshooting

### 1. Deployment Failed

- ✅ Cek semua secrets sudah di-set
- ✅ Cek GCP_SA_KEY format JSON valid
- ✅ Cek service account permissions

### 2. SMTP Not Working

- ✅ Cek SMTP_PASSWORD sudah benar
- ✅ Cek 2FA Gmail aktif
- ✅ Cek App Password valid

### 3. Database Connection Failed

- ✅ Cek MONGODB_URI format
- ✅ Cek network access
- ✅ Cek credentials valid

## Best Practices

### 1. Secret Management

- 🔒 **Never commit secrets** ke repository
- 🔒 **Use descriptive names** untuk secrets
- 🔒 **Rotate regularly** (setiap 3-6 bulan)
- 🔒 **Limit access** hanya ke yang perlu

### 2. Environment Separation

```
Development:
- SMTP_USERNAME=dev@gmail.com
- MONGODB_URI=mongodb://dev-db

Production:
- SMTP_USERNAME=prod@gmail.com
- MONGODB_URI=mongodb://prod-db
```

### 3. Monitoring

- 📊 Monitor deployment logs
- 📊 Set up alerts untuk failures
- 📊 Track secret usage

## Alternative: Manual Deployment

Jika tidak ingin menggunakan GitHub Actions, bisa deploy manual:

```bash
gcloud functions deploy sakha \
  --gen2 \
  --runtime=go121 \
  --region=asia-southeast2 \
  --source=. \
  --entry-point=URL \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars="SMTP_HOST=smtp.gmail.com,SMTP_PORT=587,SMTP_USERNAME=sakhaclothing@gmail.com,SMTP_PASSWORD=your_password,FROM_EMAIL=sakhaclothing@gmail.com,FROM_NAME=Sakha Clothing,MONGODB_URI=your_mongodb_uri,JWT_SECRET=your_jwt_secret"
```

## Checklist

- [ ] GitHub Secrets setup
- [ ] Gmail App Password dibuat
- [ ] Google Cloud Service Account dibuat
- [ ] Service Account key di-download
- [ ] Workflow file dibuat
- [ ] Test deployment
- [ ] Test email sending
- [ ] Monitor logs
