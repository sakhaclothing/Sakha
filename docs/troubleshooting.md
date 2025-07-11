# Troubleshooting Guide

## Common Issues and Solutions

### 1. Response Header Too Large Error

**Error:** `Response header too large`

**Cause:** This error occurs when the HTTP response headers exceed the maximum allowed size.

**Solutions:**

#### A. Simplified CORS Configuration
The CORS middleware has been simplified to prevent header size issues:

```go
// Minimal CORS headers
c.Set("Access-Control-Allow-Origin", "*")
```

#### B. Reduced Response Payload
- Removed `reset_link` from forgot-password response
- Simplified error messages
- Optimized Fiber configuration

#### C. Fiber Configuration
```go
app := fiber.New(fiber.Config{
    DisableStartupMessage: true,
    ServerHeader:          "Sakha API",
    AppName:              "Sakha Clothing API",
    ReadTimeout:           30,
    WriteTimeout:          30,
    IdleTimeout:           120,
    ReadBufferSize:        4096,
    WriteBufferSize:       4096,
})
```

### 2. Testing the Forgot Password Endpoint

#### Using cURL:
```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "zenkun.enterkill13@gmail.com"}' \
  -v
```

#### Using Postman:
1. Method: POST
2. URL: `https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password`
3. Headers: `Content-Type: application/json`
4. Body (raw JSON):
```json
{
  "email": "zenkun.enterkill13@gmail.com"
}
```

### 3. Expected Response

**Success Response:**
```json
{
  "message": "Link reset password telah dikirim ke email Anda"
}
```

**Error Response:**
```json
{
  "error": "Error message description"
}
```

### 4. Debugging Steps

1. **Check Logs:** The application now includes detailed logging for forgot-password requests
2. **Verify Email:** Ensure the email exists in the database
3. **Check Database Connection:** Verify MongoDB connection is working
4. **Test with Different Email:** Try with a different email address

### 5. Common HTTP Status Codes

- `200` - Success
- `400` - Bad Request (validation errors)
- `500` - Internal Server Error

### 6. Environment Variables

Make sure these are set correctly:
```bash
MONGODB_URI=your_mongodb_connection_string
JWT_SECRET=your_jwt_secret
```

### 7. Database Collections

Ensure these collections exist:
- `users` - User accounts
- `password_resets` - Password reset tokens

### 8. Email Configuration (Optional)

For production email sending:
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

### 9. Testing Checklist

- [ ] Database connection is working
- [ ] Email exists in users collection
- [ ] CORS headers are minimal
- [ ] Response payload is small
- [ ] No extra headers are being set
- [ ] Fiber configuration is optimized

### 10. Alternative Testing

If the main endpoint still fails, try testing with a simpler request:

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

This will help determine if the issue is specific to the forgot-password endpoint or a general configuration problem. 