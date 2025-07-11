# Minimal Testing Guide

## Ultra-Minimal Endpoints for Debugging

### 1. Test Basic Endpoint

**POST** `/test`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/test \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response:** `OK`

### 2. Ultra-Minimal Forgot Password

**POST** `/auth/forgot-password`

```bash
curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "zenkun.enterkill13@gmail.com"}'
```

**Expected Responses:**

- Success: `Email sent`
- Email not found: `Email sent if registered`
- Invalid data: `Invalid data`
- Database error: `Database error`
- Token error: `Token error`
- Save error: `Save error`

## Changes Made

### 1. Ultra-Minimal Fiber Config

```go
app := fiber.New(fiber.Config{
    DisableStartupMessage:     true,
    DisableDefaultDate:        true,
    DisableDefaultContentType: true,
})
```

### 2. Minimal CORS

```go
app.Use(func(c *fiber.Ctx) error {
    if c.Method() == "OPTIONS" {
        c.Set("Access-Control-Allow-Origin", "*")
        c.Set("Access-Control-Allow-Methods", "POST")
        c.Set("Access-Control-Allow-Headers", "Content-Type")
        return c.SendStatus(http.StatusNoContent)
    }
    return c.Next()
})
```

### 3. Plain Text Responses

- Using `c.SendString()` instead of `c.JSON()`
- No JSON formatting overhead
- Minimal response headers

## Testing Steps

1. **Test Basic Endpoint First**

   ```bash
   curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/test
   ```

2. **If Basic Test Works, Test Forgot Password**

   ```bash
   curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
     -H "Content-Type: application/json" \
     -d '{"email": "zenkun.enterkill13@gmail.com"}'
   ```

3. **Check Response Headers**
   ```bash
   curl -X POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
     -H "Content-Type: application/json" \
     -d '{"email": "zenkun.enterkill13@gmail.com"}' \
     -v
   ```

## Troubleshooting

### If Still Getting "Response Header Too Large"

1. **Check if it's a Google Cloud Functions limitation**

   - Try the basic `/test` endpoint first
   - If that fails, it's a platform issue

2. **Check if it's specific to forgot-password**

   - Try other endpoints like `/auth/login`
   - Compare response headers

3. **Check if it's the email parameter**
   - Try with a shorter email
   - Try with different email formats

### Alternative Testing

If the endpoint still fails, try testing with a different tool:

```bash
# Using wget
wget --post-data='{"email":"test@example.com"}' \
  --header='Content-Type: application/json' \
  https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password

# Using httpie
http POST https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/auth/forgot-password \
  email=test@example.com
```

## Expected Behavior

- **Success:** Should return plain text "Email sent"
- **No JSON:** No JSON formatting to reduce header size
- **Minimal Headers:** Only essential CORS headers
- **Fast Response:** Should be very fast due to minimal processing

## Debug Information

The endpoint now:

- Uses plain text responses
- Has minimal error handling
- Removes all logging
- Uses ultra-minimal Fiber config
- Has minimal CORS headers
- Removes JSON formatting overhead
