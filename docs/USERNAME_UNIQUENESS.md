# Username Uniqueness Implementation

## Overview

Sistem ini telah diimplementasikan dengan logika untuk memastikan bahwa setiap username hanya dapat digunakan oleh satu user saja (tidak ada duplikasi).

## Fitur yang Diimplementasikan

### 1. Database Level Uniqueness

- **Unique Index**: Dibuat unique index pada field `username` di MongoDB
- **Automatic Creation**: Index dibuat otomatis saat aplikasi start
- **Case Insensitive**: Username disimpan dalam lowercase untuk konsistensi

### 2. Application Level Validation

- **Pre-registration Check**: Validasi username sebelum registrasi
- **Case Insensitive Search**: Pencarian username tidak case sensitive
- **Format Validation**: Validasi format username (huruf, angka, underscore)
- **Length Validation**: Username harus 3-20 karakter

### 3. API Endpoints

#### Register User

```
POST /auth/register
```

**Request Body:**

```json
{
  "username": "user123",
  "password": "password123",
  "email": "user@example.com",
  "fullname": "User Name"
}
```

**Response (Success - 201):**

```json
{
  "message": "Berhasil register",
  "user_id": "507f1f77bcf86cd799439011"
}
```

**Response (Error - 409):**

```json
{
  "error": "Username sudah digunakan"
}
```

#### Check Username Availability

```
POST /auth/check-username
```

**Request Body:**

```json
{
  "username": "user123"
}
```

**Response (Available):**

```json
{
  "available": true,
  "message": "Username tersedia"
}
```

**Response (Not Available):**

```json
{
  "available": false,
  "message": "Username sudah digunakan"
}
```

#### Login

```
POST /auth/login
```

**Request Body:**

```json
{
  "username": "user123",
  "password": "password123"
}
```

## Validasi Username

### Format yang Diizinkan:

- **Karakter**: Hanya huruf (a-z), angka (0-9), dan underscore (\_)
- **Panjang**: 3-20 karakter
- **Case**: Otomatis dikonversi ke lowercase

### Contoh Username Valid:

- `user123`
- `john_doe`
- `admin2024`
- `test_user`

### Contoh Username Tidak Valid:

- `user@123` (mengandung karakter khusus)
- `ab` (terlalu pendek)
- `very_long_username_that_exceeds_limit` (terlalu panjang)
- `User Name` (mengandung spasi)

## Implementasi Teknis

### 1. Database Index

```go
// Di config/config.go
func CreateUniqueIndexes() {
    usernameIndexModel := mongo.IndexModel{
        Keys: map[string]interface{}{
            "username": 1,
        },
        Options: options.Index().SetUnique(true).SetName("username_unique"),
    }

    _, err := DB.Collection("users").Indexes().CreateOne(ctx, usernameIndexModel)
}
```

### 2. Username Validation

```go
// Di controller/auth_controller.go
func isValidUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }

    for _, char := range username {
        if !((char >= 'a' && char <= 'z') ||
             (char >= '0' && char <= '9') ||
             char == '_') {
            return false
        }
    }

    return true
}
```

### 3. Case Insensitive Search

```go
// Pencarian username dengan regex case insensitive
err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
    "username": bson.M{"$regex": "^" + username + "$", "$options": "i"},
}).Decode(&existingUser)
```

## Testing

### Menjalankan Test

```bash
cd tests
go run test_username_uniqueness.go
```

### Test Cases yang Dicakup:

1. ✅ Check username availability
2. ✅ Register first user
3. ✅ Check same username again (should be unavailable)
4. ✅ Try to register with same username (should fail)
5. ✅ Try to register with same username but different case (should fail)
6. ✅ Register with different username (should succeed)
7. ✅ Test invalid username format

## Error Handling

### HTTP Status Codes:

- **200**: Success (check username, login)
- **201**: Created (register success)
- **400**: Bad Request (invalid input)
- **401**: Unauthorized (wrong password)
- **404**: Not Found (user not found)
- **409**: Conflict (username already exists)
- **500**: Internal Server Error

### Error Messages:

- `"Username sudah digunakan"` - Username sudah ada
- `"Username hanya boleh berisi huruf, angka, dan underscore"` - Format tidak valid
- `"Username tidak boleh kosong"` - Username kosong
- `"Username atau password salah"` - Login gagal

## Security Considerations

1. **Input Sanitization**: Username di-trim dan dikonversi ke lowercase
2. **SQL Injection Prevention**: Menggunakan MongoDB driver yang aman
3. **Rate Limiting**: Disarankan untuk menambahkan rate limiting pada endpoint register
4. **Password Security**: Password di-hash menggunakan bcrypt

## Monitoring dan Logging

Sistem mencatat:

- ✅ Koneksi database berhasil
- ✅ Unique index berhasil dibuat
- ⚠️ Error saat membuat index (normal jika sudah ada)
- ❌ Error database operations

## Deployment Notes

1. **First Run**: Index akan dibuat otomatis saat aplikasi pertama kali dijalankan
2. **Existing Data**: Jika sudah ada data duplikat, index creation akan gagal
3. **Migration**: Pastikan tidak ada username duplikat sebelum deploy
4. **Backup**: Backup database sebelum implementasi

## Troubleshooting

### Index Creation Failed

```
⚠️ Gagal membuat unique index untuk username: E11000 duplicate key error
```

**Solution**: Hapus username duplikat dari database terlebih dahulu

### Username Not Found During Login

```
❌ Username atau password salah
```

**Solution**: Pastikan username yang dimasukkan sudah benar (case insensitive)

### Registration Always Fails

```
❌ Username sudah digunakan
```

**Solution**: Cek apakah ada user dengan username yang sama (termasuk case yang berbeda)
