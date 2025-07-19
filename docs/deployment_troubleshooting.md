# Deployment Troubleshooting - Sakha Clothing

## Common Deployment Issues

### 1. Import Path Errors

#### Error: `package Sakha/config is not in std`

```
ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Build failed with status: FAILURE and message: controller/newsletter_controller.go:14:2: package Sakha/config is not in std
```

#### Solution:

**Problem**: Import path menggunakan `"Sakha/config"` instead of `"github.com/sakhaclothing/config"`

**Fix**: Update import paths in all files to use correct module path:

```go
// ❌ Wrong
import (
    "Sakha/config"
    "Sakha/model"
    "Sakha/utils"
)

// ✅ Correct
import (
    "github.com/sakhaclothing/config"
    "github.com/sakhaclothing/model"
    "github.com/sakhaclothing/utils"
)
```

**Files to check**:

- `Sakha/controller/newsletter_controller.go`
- `Sakha/controller/product_controller.go`
- `Sakha/route/route.go`
- `Sakha/utils/newsletter.go`

### 2. Undefined Function Errors

#### Error: `undefined: sendNewProductNotificationToSubscribers`

```
controller/product_controller.go:157: undefined: sendNewProductNotificationToSubscribers
```

#### Solution:

**Problem**: Function called but not defined in the same package

**Fix**: Move functions to utils package and import them:

```go
// In product_controller.go
import (
    "github.com/sakhaclothing/utils"
)

// Call function
go utils.SendNewProductNotificationToSubscribers(product)
```

### 3. Module Path Issues

#### Error: `go.mod` not found or incorrect

```
go: cannot find main module
```

#### Solution:

**Check `go.mod` file**:

```go
module github.com/sakhaclothing

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.52.0
    go.mongodb.org/mongo-driver v1.13.1
    // other dependencies...
)
```

### 4. Missing Dependencies

#### Error: `cannot find package`

```
cannot find package "github.com/gofiber/fiber/v2"
```

#### Solution:

```bash
# Update dependencies
go mod tidy
go mod download

# Check go.mod and go.sum files
cat go.mod
cat go.sum
```

## Deployment Checklist

### ✅ Pre-Deployment Checks

1. **Import Paths**

   - [ ] All imports use `github.com/sakhaclothing/...`
   - [ ] No local `Sakha/...` imports

2. **Function Definitions**

   - [ ] All functions are defined in correct packages
   - [ ] No undefined function calls

3. **Dependencies**

   - [ ] `go.mod` exists and is correct
   - [ ] `go.sum` is up to date
   - [ ] Run `go mod tidy`

4. **File Structure**
   - [ ] All required files are present
   - [ ] No syntax errors

### ✅ Deployment Commands

```bash
# 1. Navigate to Sakha directory
cd Sakha

# 2. Update dependencies
go mod tidy
go mod download

# 3. Test build locally
go build -o main .

# 4. Deploy to Google Cloud Functions
gcloud functions deploy sakha \
  --runtime go121 \
  --trigger-http \
  --allow-unauthenticated \
  --region asia-southeast2 \
  --entry-point URL
```

## File Structure Requirements

### Required Files for Deployment

```
Sakha/
├── go.mod                    # Module definition
├── go.sum                    # Dependency checksums
├── main.go                   # Entry point
├── config/
│   ├── config.go            # Database config
│   └── email_config.go      # Email config
├── controller/
│   ├── auth_controller.go   # Auth handlers
│   ├── product_controller.go # Product handlers
│   └── newsletter_controller.go # Newsletter handlers
├── model/
│   ├── user.go              # User model
│   ├── product.go           # Product model
│   └── newsletter.go        # Newsletter models
├── route/
│   └── route.go             # Route definitions
├── utils/
│   ├── email.go             # Email utilities
│   ├── token.go             # Token utilities
│   └── newsletter.go        # Newsletter utilities
└── middlewares/
    └── jwt.go               # JWT middleware
```

## Common Fixes

### 1. Fix Import Paths

**Find and replace in all `.go` files**:

```bash
# Replace local imports with module imports
sed -i 's|"Sakha/|"github.com/sakhaclothing/|g' *.go
sed -i 's|"Sakha/|"github.com/sakhaclothing/|g' */*.go
```

### 2. Move Functions to Utils

**If functions are undefined, move them to utils package**:

```go
// Create utils/newsletter.go
package utils

import (
    "github.com/sakhaclothing/config"
    "github.com/sakhaclothing/model"
)

func SendNewProductNotificationToSubscribers(product model.Product) {
    // Implementation here
}
```

### 3. Update Function Calls

**Update all function calls to use utils package**:

```go
// In controllers
import "github.com/sakhaclothing/utils"

// Call functions
go utils.SendNewProductNotificationToSubscribers(product)
```

## Testing Before Deployment

### 1. Local Build Test

```bash
cd Sakha
go build -o main .
```

### 2. Syntax Check

```bash
go vet ./...
```

### 3. Dependency Check

```bash
go mod verify
```

### 4. Test Run

```bash
# Test locally (if possible)
go run main.go
```

## Deployment Commands

### Google Cloud Functions

```bash
# Deploy function
gcloud functions deploy sakha \
  --runtime go121 \
  --trigger-http \
  --allow-unauthenticated \
  --region asia-southeast2 \
  --entry-point URL \
  --source .

# Check deployment status
gcloud functions describe sakha --region asia-southeast2

# View logs
gcloud functions logs read sakha --region asia-southeast2
```

### Alternative: Manual Upload

```bash
# Create deployment package
zip -r sakha.zip . -x "*.git*" "*.DS_Store*" "node_modules/*"

# Upload to Google Cloud Storage
gsutil cp sakha.zip gs://your-bucket/

# Deploy from GCS
gcloud functions deploy sakha \
  --runtime go121 \
  --trigger-http \
  --allow-unauthenticated \
  --region asia-southeast2 \
  --entry-point URL \
  --source gs://your-bucket/sakha.zip
```

## Monitoring and Debugging

### 1. Check Function Logs

```bash
gcloud functions logs read sakha --region asia-southeast2 --limit 50
```

### 2. Test Function Endpoints

```bash
# Test health check
curl https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/test

# Test products endpoint
curl https://asia-southeast2-ornate-course-437014-u9.cloudfunctions.net/sakha/products
```

### 3. Check Function Status

```bash
gcloud functions describe sakha --region asia-southeast2
```

## Rollback Strategy

### 1. Keep Previous Version

```bash
# Deploy with different name
gcloud functions deploy sakha-v2 \
  --runtime go121 \
  --trigger-http \
  --allow-unauthenticated \
  --region asia-southeast2 \
  --entry-point URL
```

### 2. Rollback to Previous Version

```bash
# Delete current version
gcloud functions delete sakha --region asia-southeast2

# Deploy previous version
gcloud functions deploy sakha \
  --runtime go121 \
  --trigger-http \
  --allow-unauthenticated \
  --region asia-southeast2 \
  --entry-point URL \
  --source gs://your-bucket/sakha-previous.zip
```

## Prevention Tips

### 1. Use Consistent Import Paths

- Always use `github.com/sakhaclothing/...` imports
- Never use local `Sakha/...` imports

### 2. Organize Functions Properly

- Put utility functions in `utils/` package
- Keep controllers focused on HTTP handling
- Use proper package structure

### 3. Test Before Deploy

- Always run `go build` locally
- Check for syntax errors
- Verify dependencies

### 4. Use Version Control

- Commit working versions
- Tag releases
- Keep deployment history

## Emergency Contacts

### Google Cloud Support

- [Google Cloud Console](https://console.cloud.google.com)
- [Cloud Functions Documentation](https://cloud.google.com/functions/docs)
- [Go Runtime Documentation](https://cloud.google.com/functions/docs/writing/go)

### Project Structure

- Backend: `Sakha/` directory
- Frontend: Various directories in root
- Documentation: `Sakha/docs/`
