# Sakha API Endpoints

## Base URL

```
https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha
```

## Authentication Endpoints

### 1. Register User

**POST** `/auth/register`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com",
    "fullname": "Test User"
  }'
```

### 2. Login

**POST** `/auth/login`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. Check Username Availability

**POST** `/auth/check-username`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/check-username \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser"
  }'
```

### 4. Get User Profile

**POST** `/auth/profile`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Forgot Password Endpoints

### 5. Request Password Reset

**POST** `/auth/forgot-password`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

**Response:**

```json
{
  "message": "Link reset password telah dikirim ke email Anda",
  "reset_link": "https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/reset-password?token=abc123..."
}
```

### 6. Reset Password

**POST** `/auth/reset-password`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "abc123...",
    "password": "newpassword123"
  }'
```

**Response:**

```json
{
  "message": "Password berhasil direset"
}
```

## Testing with Postman

### Environment Variables

Set these in your Postman environment:

```
BASE_URL: https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha
```

### Collection Variables

After login, save the JWT token:

```
TOKEN: {{response.body.token}}
```

### Example Requests

1. **Register User**

   - Method: POST
   - URL: {{BASE_URL}}/auth/register
   - Body: Raw JSON

   ```json
   {
     "username": "testuser",
     "password": "password123",
     "email": "test@example.com",
     "fullname": "Test User"
   }
   ```

2. **Login**

   - Method: POST
   - URL: {{BASE_URL}}/auth/login
   - Body: Raw JSON

   ```json
   {
     "username": "testuser",
     "password": "password123"
   }
   ```

3. **Get Profile**

   - Method: POST
   - URL: {{BASE_URL}}/auth/profile
   - Headers: Authorization: Bearer {{TOKEN}}

4. **Forgot Password**

   - Method: POST
   - URL: {{BASE_URL}}/auth/forgot-password
   - Body: Raw JSON

   ```json
   {
     "email": "test@example.com"
   }
   ```

5. **Reset Password**
   - Method: POST
   - URL: {{BASE_URL}}/auth/reset-password
   - Body: Raw JSON
   ```json
   {
     "token": "token_from_email",
     "password": "newpassword123"
   }
   ```

## Error Responses

All endpoints return consistent error formats:

```json
{
  "error": "Error message description"
}
```

Common HTTP Status Codes:

- `200` - Success
- `201` - Created (for register)
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid token)
- `404` - Not Found
- `409` - Conflict (duplicate username/email)
- `500` - Internal Server Error

## Notes

1. All endpoints accept and return JSON
2. JWT tokens expire in 24 hours
3. Password reset tokens expire in 1 hour
4. Email validation is case-insensitive
5. Username validation is case-insensitive
6. Password minimum length is 6 characters
