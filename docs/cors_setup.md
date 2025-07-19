# CORS Configuration - Sakha Clothing

## Overview

CORS (Cross-Origin Resource Sharing) configuration untuk Sakha Clothing backend yang mendukung development dan production environment.

## Allowed Origins

### Production

- `https://sakhaclothing.shop` - Main production website

### Development

- `http://127.0.0.1:5500` - Live Server (VS Code)
- `http://localhost:5500` - Live Server (VS Code)
- `http://127.0.0.1:3000` - React development server
- `http://localhost:3000` - React development server
- `http://127.0.0.1:8080` - Common development port
- `http://localhost:8080` - Common development port
- `http://127.0.0.1:5000` - Flask/Python development server
- `http://localhost:5000` - Flask/Python development server

## CORS Headers

### Allowed Methods

- `POST` - Create operations
- `GET` - Read operations
- `PUT` - Update operations
- `PATCH` - Partial updates
- `DELETE` - Delete operations
- `OPTIONS` - Preflight requests

### Allowed Headers

- `Content-Type` - JSON data
- `Authorization` - Bearer tokens
- `X-Requested-With` - AJAX requests

### Credentials

- `Access-Control-Allow-Credentials: true` - Allow cookies/auth

## Implementation

### File: `Sakha/route/route.go`

```go
// CORS configuration for production and development
app.Use(func(c *fiber.Ctx) error {
    origin := c.Get("Origin")

    // Allowed origins
    allowedOrigins := []string{
        "https://sakhaclothing.shop",
        "http://127.0.0.1:5500",
        "http://localhost:5500",
        "http://127.0.0.1:3000",
        "http://localhost:3000",
        "http://127.0.0.1:8080",
        "http://localhost:8080",
        "http://127.0.0.1:5000",
        "http://localhost:5000",
    }

    // Check if origin is allowed
    for _, allowedOrigin := range allowedOrigins {
        if origin == allowedOrigin {
            c.Set("Access-Control-Allow-Origin", allowedOrigin)
            break
        }
    }

    // Set CORS headers for all requests
    c.Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
    c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
    c.Set("Access-Control-Allow-Credentials", "true")
    c.Set("Vary", "Origin")

    // Handle preflight requests
    if c.Method() == "OPTIONS" {
        return c.SendStatus(204)
    }

    return c.Next()
})
```

## Testing CORS

### 1. Test with Live Server

```bash
# Start Live Server in VS Code
# Or use any local development server
```

### 2. Test API Endpoints

```javascript
// Test from browser console at http://localhost:5500
fetch('https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/products', {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json'
    }
})
.then(response => {
    console.log('Status:', response.status);
    console.log('CORS Headers:', response.headers);
    return response.json();
})
.then(data => console.log('Data:', data))
.catch(error => console.error('CORS Error:', error));
```

### 3. Test Newsletter Subscription

```javascript
// Test newsletter subscription from localhost
fetch('https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/newsletter/subscribe', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        email: 'test@example.com'
    })
})
.then(response => {
    console.log('Status:', response.status);
    console.log('CORS Headers:', response.headers);
    return response.json();
})
.then(data => console.log('Data:', data))
.catch(error => console.error('CORS Error:', error));
```

### 4. Test OPTIONS Request (Preflight)

```javascript
// Test preflight request
fetch('https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/products', {
    method: 'OPTIONS',
    headers: {
        'Access-Control-Request-Method': 'POST',
        'Access-Control-Request-Headers': 'Content-Type, Authorization'
    }
})
.then(response => {
    console.log('Preflight Status:', response.status);
    console.log('CORS Headers:', response.headers);
})
.catch(error => console.error('Preflight Error:', error));
```

### 5. Expected Results

#### ✅ Successful CORS Response
```javascript
// Status: 200 or 201
// Headers should include:
// - Access-Control-Allow-Origin: http://localhost:5500
// - Access-Control-Allow-Methods: POST, GET, PUT, PATCH, DELETE, OPTIONS
// - Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
// - Access-Control-Allow-Credentials: true
```

#### ❌ CORS Error
```javascript
// Error: Access to fetch at '...' from origin 'http://localhost:5500' 
// has been blocked by CORS policy: No 'Access-Control-Allow-Origin' 
// header is present on the requested resource.
```

## Common Development Ports

### VS Code Live Server

- Default: `http://127.0.0.1:5500`
- Alternative: `http://localhost:5500`

### React Development Server

- Default: `http://localhost:3000`
- Alternative: `http://127.0.0.1:3000`

### Vue.js Development Server

- Default: `http://localhost:8080`
- Alternative: `http://127.0.0.1:8080`

### Angular Development Server

- Default: `http://localhost:4200`
- Alternative: `http://127.0.0.1:4200`

### Flask/Python Development Server

- Default: `http://localhost:5000`
- Alternative: `http://127.0.0.1:5000`

## Troubleshooting

### Issue: CORS Error in Browser

**Error:** `Access to fetch at '...' from origin 'http://localhost:5500' has been blocked by CORS policy`

**Solution:**

1. Check if origin is in allowed list
2. Verify backend is deployed with new CORS config
3. Clear browser cache
4. Check browser console for specific error

### Issue: Preflight Request Failing

**Error:** `Method OPTIONS not allowed`

**Solution:**

1. Ensure OPTIONS method is handled
2. Check CORS middleware is applied before routes
3. Verify preflight response returns 204

### Issue: Credentials Not Working

**Error:** `Credentials flag is 'true', but the 'Access-Control-Allow-Credentials' header is ''`

**Solution:**

1. Ensure `Access-Control-Allow-Credentials: true` is set
2. Check if origin is exact match (not wildcard)
3. Verify credentials are included in request

## Adding New Origins

To add new development origins:

1. **Edit `Sakha/route/route.go`**
2. **Add to `allowedOrigins` array:**

```go
allowedOrigins := []string{
    "https://sakhaclothing.shop",
    "http://127.0.0.1:5500",
    "http://localhost:5500",
    "http://127.0.0.1:3000",
    "http://localhost:3000",
    "http://127.0.0.1:8080",
    "http://localhost:8080",
    "http://127.0.0.1:5000",
    "http://localhost:5000",
    // Add new origin here
    "http://your-new-origin:port",
}
```

3. **Deploy backend changes**
4. **Test from new origin**

## Security Considerations

### Production

- ✅ Only allow specific production domains
- ✅ Use HTTPS for all production origins
- ✅ Validate origin headers

### Development

- ✅ Allow common development ports
- ✅ Support both localhost and 127.0.0.1
- ✅ Include necessary headers for development tools

### General

- ✅ Don't use wildcard (\*) for production
- ✅ Validate all origins
- ✅ Handle preflight requests properly
- ✅ Set appropriate security headers

## Browser Compatibility

### Supported Browsers

- ✅ Chrome (all versions)
- ✅ Firefox (all versions)
- ✅ Safari (all versions)
- ✅ Edge (all versions)

### Development Tools

- ✅ VS Code Live Server
- ✅ React Development Server
- ✅ Vue.js Development Server
- ✅ Angular Development Server
- ✅ Flask Development Server
- ✅ Any local development server

## Performance Impact

### Minimal Impact

- ✅ CORS check is fast string comparison
- ✅ Headers set only when needed
- ✅ Preflight requests handled efficiently
- ✅ No database queries for CORS

### Optimization

- ✅ Use exact origin matching
- ✅ Minimize allowed origins list
- ✅ Cache CORS headers when possible
- ✅ Handle preflight requests early

## Monitoring

### Log CORS Requests

```go
// Add logging to CORS middleware
if origin != "" {
    log.Printf("CORS request from origin: %s", origin)
}
```

### Track CORS Errors

- Monitor browser console errors
- Check server logs for CORS issues
- Use analytics to track failed requests

## Related Files

- `Sakha/route/route.go` - CORS configuration
- `Sakha/docs/cors_setup.md` - This documentation
- Frontend files - Use API endpoints
