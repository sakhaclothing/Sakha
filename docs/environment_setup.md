# Environment Setup Guide

## Overview

Dokumen ini menjelaskan cara setup environment variables untuk project Sakha Clothing.

## Environment Variables

### Required Variables

#### MONGOSTRING

Connection string untuk MongoDB database.

**Format:**

```
MONGOSTRING=mongodb://host:port/database
```

**Contoh:**

**MongoDB Lokal:**

```
MONGOSTRING=mongodb://localhost:27017/sakha
```

**MongoDB Atlas:**

```
MONGOSTRING=mongodb+srv://username:password@cluster.mongodb.net/sakha?retryWrites=true&w=majority
```

**MongoDB dengan Authentication:**

```
MONGOSTRING=mongodb://username:password@localhost:27017/sakha
```

### Optional Variables

#### ENV

Environment mode (development/production)

```
ENV=development
```

#### PORT

Server port (default: 8080)

```
PORT=8080
```

## Setup Methods

### Method 1: .env File (Recommended for Development)

1. **Buat file `.env` di root directory Sakha:**

```bash
cd Sakha
touch .env
```

2. **Isi file `.env` dengan konfigurasi:**

```env
# MongoDB Connection String
MONGOSTRING=mongodb://localhost:27017/sakha

# Environment
ENV=development

# Server Port
PORT=8080
```

3. **File `.env` sudah di-ignore oleh git untuk keamanan**

### Method 2: System Environment Variables

#### Windows (PowerShell)

```powershell
$env:MONGOSTRING="mongodb://localhost:27017/sakha"
$env:ENV="development"
$env:PORT="8080"
```

#### Windows (Command Prompt)

```cmd
set MONGOSTRING=mongodb://localhost:27017/sakha
set ENV=development
set PORT=8080
```

#### Linux/macOS

```bash
export MONGOSTRING="mongodb://localhost:27017/sakha"
export ENV="development"
export PORT="8080"
```

### Method 3: Docker Environment

Jika menggunakan Docker, tambahkan ke `docker-compose.yml`:

```yaml
environment:
  - MONGOSTRING=mongodb://mongo:27017/sakha
  - ENV=development
  - PORT=8080
```

### Method 4: Cloud Functions (Google Cloud)

Jika deploy ke Google Cloud Functions, set environment variables di console:

1. Buka Google Cloud Console
2. Pilih Cloud Functions
3. Edit function
4. Set environment variables:
   - `MONGOSTRING`: Your MongoDB connection string
   - `ENV`: production
   - `PORT`: 8080

## MongoDB Setup

### Local MongoDB

1. **Install MongoDB:**

```bash
# Ubuntu/Debian
sudo apt-get install mongodb

# macOS dengan Homebrew
brew install mongodb-community

# Windows
# Download dari https://www.mongodb.com/try/download/community
```

2. **Start MongoDB service:**

```bash
# Ubuntu/Debian
sudo systemctl start mongodb

# macOS
brew services start mongodb-community

# Windows
# Start MongoDB service dari Services
```

3. **Test connection:**

```bash
mongosh "mongodb://localhost:27017/sakha"
```

### MongoDB Atlas (Cloud)

1. **Buat account di MongoDB Atlas**
2. **Buat cluster baru**
3. **Set up database access:**
   - Username dan password
4. **Set up network access:**
   - Allow access from anywhere (0.0.0.0/0) untuk development
5. **Get connection string:**
   - Klik "Connect"
   - Pilih "Connect your application"
   - Copy connection string

## Testing Connection

### Test dengan Go

```go
package main

import (
    "log"
    "github.com/sakhaclothing/config"
)

func main() {
    config.ConnectDB()
    log.Println("✅ Database connection successful!")
}
```

### Test dengan MongoDB Shell

```bash
# Test local connection
mongosh "mongodb://localhost:27017/sakha"

# Test Atlas connection
mongosh "mongodb+srv://username:password@cluster.mongodb.net/sakha"
```

## Troubleshooting

### Error: "ENV MONGOSTRING tidak ditemukan"

**Solution:**

1. Pastikan file `.env` ada di root directory Sakha
2. Pastikan format file `.env` benar
3. Restart aplikasi setelah mengubah environment variables

### Error: "Failed to connect to MongoDB"

**Solutions:**

1. **Check MongoDB service:**

```bash
# Ubuntu/Debian
sudo systemctl status mongodb

# macOS
brew services list | grep mongodb
```

2. **Check connection string:**

   - Pastikan host dan port benar
   - Pastikan username dan password benar (jika ada)
   - Pastikan database name benar

3. **Check network:**
   - Pastikan MongoDB accessible dari aplikasi
   - Check firewall settings
   - Check MongoDB bind IP

### Error: "Authentication failed"

**Solutions:**

1. **Check credentials:**

   - Pastikan username dan password benar
   - Pastikan user memiliki akses ke database

2. **Check authentication database:**

```bash
# Login dengan authentication database yang benar
mongosh "mongodb://username:password@host:port/authDatabase"
```

## Security Best Practices

### Development

- Gunakan MongoDB lokal tanpa authentication
- File `.env` tidak di-commit ke git
- Gunakan connection string sederhana

### Production

- Gunakan MongoDB Atlas atau managed service
- Set up proper authentication
- Use connection string dengan username/password
- Restrict network access
- Use environment variables di cloud platform

## Example Configurations

### Development (.env)

```env
MONGOSTRING=mongodb://localhost:27017/sakha
ENV=development
PORT=8080
```

### Production (Environment Variables)

```bash
MONGOSTRING=mongodb+srv://produser:prodpass@cluster.mongodb.net/sakha_prod
ENV=production
PORT=8080
```

### Docker (.env)

```env
MONGOSTRING=mongodb://mongo:27017/sakha
ENV=development
PORT=8080
```

## Next Steps

1. **Setup MongoDB** (local atau cloud)
2. **Buat file `.env`** dengan connection string
3. **Test connection** dengan aplikasi
4. **Deploy** dengan environment variables yang sesuai
