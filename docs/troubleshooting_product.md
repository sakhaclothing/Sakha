# Troubleshooting Product Management

## Common Issues and Solutions

### 1. Cannot Edit Products

#### Problem

- Edit button doesn't work
- Modal doesn't open
- Product not found error

#### Solutions

**Check Browser Console:**

1. Open browser developer tools (F12)
2. Go to Console tab
3. Look for error messages when clicking edit button
4. Check if `productManager` is defined

**Check Product ID:**

```javascript
// In browser console, check if products are loaded
console.log(productManager.products);
```

**Verify API Response:**

1. Open Network tab in developer tools
2. Click edit button
3. Check if API call is made
4. Verify response format

**Common Issues:**

- Product ID mismatch (MongoDB uses `_id`, frontend might expect `id`)
- API endpoint not found (404 error)
- CORS issues
- Network connectivity problems

### 2. Cannot Delete Products

#### Problem

- Delete button doesn't work
- Confirmation dialog doesn't appear
- Product not deleted after confirmation

#### Solutions

**Check SweetAlert2:**

```html
<!-- Make sure SweetAlert2 is loaded -->
<script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
```

**Check API Endpoint:**

```javascript
// Test delete endpoint manually
fetch("https://your-api-url/products/PRODUCT_ID", {
  method: "DELETE",
})
  .then((response) => response.json())
  .then((data) => console.log(data));
```

**Verify Product ID:**

```javascript
// Check if product ID is correct
console.log("Product ID:", productId);
```

### 3. Products Not Loading

#### Problem

- Empty product list
- Loading error
- API connection failed

#### Solutions

**Check API URL:**

```javascript
// Verify API base URL is correct
console.log("API Base URL:", productManager.apiBaseUrl);
```

**Test API Endpoint:**

```bash
# Test with curl
curl -X GET https://your-api-url/products
```

**Check Network Tab:**

1. Open developer tools
2. Go to Network tab
3. Refresh page
4. Look for failed requests

**Common Issues:**

- Wrong API URL
- Backend server not running
- CORS configuration
- Database connection issues

### 4. Featured Products Not Showing

#### Problem

- Featured products page shows empty
- Fallback products showing instead

#### Solutions

**Check Featured Query:**

```javascript
// Verify featured query parameter
const response = await fetch(`${apiBaseUrl}/products?featured=true`);
```

**Check Database:**

```javascript
// Verify products have is_featured: true
db.products.find({ is_featured: true, is_active: true });
```

**Check API Response:**

```javascript
// Test featured products endpoint
fetch("https://your-api-url/products?featured=true")
  .then((response) => response.json())
  .then((data) => console.log(data));
```

### 5. Form Validation Issues

#### Problem

- Form submission fails
- Validation errors
- Required fields not working

#### Solutions

**Check Form Fields:**

```html
<!-- Ensure required fields have required attribute -->
<input type="text" id="productName" required />
<input type="number" id="productPrice" required min="0" />
```

**Check Form Data:**

```javascript
// Log form data before submission
console.log("Form data:", formData);
```

**Validate Backend:**

```go
// Check backend validation
if strings.TrimSpace(product.Name) == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "status":  "error",
        "message": "Product name is required",
    })
}
```

## Debug Steps

### Step 1: Check Browser Console

1. Open developer tools (F12)
2. Go to Console tab
3. Look for JavaScript errors
4. Check if all scripts are loaded

### Step 2: Check Network Requests

1. Go to Network tab
2. Perform the action (edit/delete)
3. Look for failed requests
4. Check response status and data

### Step 3: Test API Endpoints

Use the test page: `/dashboard/test-api.html`

### Step 4: Check Database

```bash
# Connect to MongoDB
mongosh "your-connection-string"

# Check products collection
use sakha
db.products.find().pretty()
```

### Step 5: Verify Backend Logs

Check your backend server logs for errors.

## Common Error Messages

### "Product not found"

- Product ID is incorrect
- Product was deleted
- Database connection issue

### "Invalid product ID"

- ID format is wrong
- Not a valid MongoDB ObjectId

### "Failed to fetch products"

- Database connection failed
- Collection doesn't exist
- Permission issues

### "CORS error"

- Backend CORS configuration
- Wrong API URL
- Protocol mismatch (http vs https)

## Testing Checklist

### Frontend Testing

- [ ] Page loads without errors
- [ ] Products are displayed
- [ ] Add product works
- [ ] Edit product works
- [ ] Delete product works
- [ ] Search works
- [ ] Filter works
- [ ] Featured toggle works

### Backend Testing

- [ ] API endpoints respond
- [ ] Database connection works
- [ ] CRUD operations work
- [ ] Validation works
- [ ] Error handling works

### Integration Testing

- [ ] Frontend can communicate with backend
- [ ] Data flows correctly
- [ ] Real-time updates work
- [ ] Error messages are displayed

## Performance Issues

### Slow Loading

- Check database indexes
- Implement pagination
- Optimize queries
- Use caching

### Memory Issues

- Check for memory leaks
- Optimize JavaScript
- Reduce DOM manipulation

## Security Issues

### Authentication

- Ensure admin authentication
- Check authorization
- Validate user permissions

### Input Validation

- Sanitize user input
- Validate on both frontend and backend
- Prevent SQL injection (MongoDB injection)

### CORS

- Configure CORS properly
- Restrict allowed origins
- Handle preflight requests

## Getting Help

If you're still having issues:

1. **Check the logs** - Both frontend and backend
2. **Use the test page** - `/dashboard/test-api.html`
3. **Verify configuration** - API URLs, database connection
4. **Check dependencies** - Make sure all packages are installed
5. **Test in isolation** - Test each component separately

## Emergency Fixes

### Quick API Test

```bash
# Test if API is working
curl -X GET https://your-api-url/products
```

### Reset Database

```bash
# Clear products collection (be careful!)
db.products.drop()
```

### Fallback Mode

If API is down, the frontend will show fallback products automatically.
