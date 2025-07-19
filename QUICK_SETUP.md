# Quick Setup Guide - Sakha Clothing

## Problem: "ENV MONGOSTRING tidak ditemukan"

### Solution 1: Create .env File (Recommended)

1. **Buat file `.env` di folder Sakha:**

```bash
cd Sakha
```

2. **Buat file `.env` dengan isi:**

```env
MONGOSTRING=mongodb://localhost:27017/sakha
ENV=development
PORT=8080
```

3. **Restart aplikasi**

### Solution 2: Set Environment Variable

#### Windows (PowerShell):

```powershell
$env:MONGOSTRING="mongodb://localhost:27017/sakha"
```

#### Windows (Command Prompt):

```cmd
set MONGOSTRING=mongodb://localhost:27017/sakha
```

#### Linux/macOS:

```bash
export MONGOSTRING="mongodb://localhost:27017/sakha"
```

### Solution 3: Install MongoDB (if not installed)

#### Windows:

1. Download dari: https://www.mongodb.com/try/download/community
2. Install MongoDB
3. Start MongoDB service

#### macOS:

```bash
brew install mongodb-community
brew services start mongodb-community
```

#### Ubuntu/Debian:

```bash
sudo apt-get install mongodb
sudo systemctl start mongodb
```

### Solution 4: Use MongoDB Atlas (Cloud)

1. Buat account di [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Buat cluster baru
3. Get connection string
4. Set environment variable:

```bash
export MONGOSTRING="mongodb+srv://username:password@cluster.mongodb.net/sakha"
```

## Test Database Connection

### Run Test Script:

```bash
cd Sakha/scripts
go run test_db_connection.go
```

### Expected Output:

```
🔍 Testing database connection...
⚠️  ENV MONGOSTRING tidak ditemukan, menggunakan MongoDB lokal
📝 Menggunakan connection string: mongodb://localhost:27017/sakha
💡 Untuk production, set environment variable MONGOSTRING
✅ Terhubung ke MongoDB
📋 Testing basic database operations...
✅ Ping successful
✅ Collections found: []
📦 Testing products collection...
✅ Products count: 0
✅ Database connection test completed!
```

## Seed Sample Data

### Run Seeding Script:

```bash
cd Sakha/scripts
go run seed_products.go
```

### Expected Output:

```
✅ Terhubung ke MongoDB
Successfully inserted product: KAOS SABLON PREMIUM (ID: ...)
Successfully inserted product: JAKET SABLON PREMIUM (ID: ...)
Successfully inserted product: CUSTOM SWEATER SABLON (ID: ...)
Successfully inserted product: CUSTOM SPORT SABLON (ID: ...)
Successfully inserted product: KAOS POLO PREMIUM (ID: ...)
Product seeding completed!
```

## Test API Endpoints

### Use Test Page:

1. Buka `/dashboard/test-api.html`
2. Klik "Test GET /products"
3. Klik "Test POST /products"
4. Klik "Test GET /products?featured=true"

### Expected Results:

- ✅ All tests should pass
- ✅ Products should be returned
- ✅ Featured products should be filtered

## Common Issues

### Issue 1: MongoDB not running

**Error:** `Failed to connect to MongoDB`

**Solution:**

```bash
# Check if MongoDB is running
sudo systemctl status mongodb  # Linux
brew services list | grep mongodb  # macOS
# Windows: Check Services app
```

### Issue 2: Wrong connection string

**Error:** `Invalid connection string`

**Solution:**

- Check format: `mongodb://host:port/database`
- For Atlas: `mongodb+srv://username:password@cluster.mongodb.net/database`

### Issue 3: Permission denied

**Error:** `Authentication failed`

**Solution:**

- Check username/password
- Ensure user has access to database
- Check authentication database

## Next Steps

1. ✅ **Setup MongoDB** (local atau cloud)
2. ✅ **Create .env file** dengan connection string
3. ✅ **Test connection** dengan script
4. ✅ **Seed sample data** (optional)
5. ✅ **Test API endpoints**
6. ✅ **Start using product management**

## Support

Jika masih ada masalah:

1. Check logs di console
2. Run test scripts
3. Verify MongoDB connection
4. Check environment variables
5. Review troubleshooting guide: `docs/troubleshooting_product.md`
