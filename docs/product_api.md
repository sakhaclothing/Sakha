# Product Management API Documentation

## Base URL

```
https://your-backend-url.com
```

## Authentication

All endpoints require authentication. Include your authentication token in the request headers.

## Endpoints

### 1. Get All Products

**GET** `/products`

Get products based on query parameters.

**Query Parameters:**

- `featured` (optional): Set to `true` to get only featured and active products
- `all` (optional): Set to `true` to get all products (active and inactive) - for admin dashboard
- If no parameters: Get only active products

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "KAOS SABLON PREMIUM",
      "description": "Kaos berkualitas tinggi dengan sablon premium",
      "price": 75000,
      "image_url": "https://example.com/image.jpg",
      "category": "kaos",
      "stock": 50,
      "is_featured": true,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 2. Get Product by ID

**GET** `/products/:id`

Get a specific product by its ID.

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "KAOS SABLON PREMIUM",
    "description": "Kaos berkualitas tinggi dengan sablon premium",
    "price": 75000,
    "image_url": "https://example.com/image.jpg",
    "category": "kaos",
    "stock": 50,
    "is_featured": true,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. Create Product

**POST** `/products`

Create a new product.

**Request Body:**

```json
{
  "name": "KAOS SABLON PREMIUM",
  "description": "Kaos berkualitas tinggi dengan sablon premium",
  "price": 75000,
  "image_url": "https://example.com/image.jpg",
  "category": "kaos",
  "stock": 50,
  "is_featured": false,
  "is_active": true
}
```

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "KAOS SABLON PREMIUM",
    "description": "Kaos berkualitas tinggi dengan sablon premium",
    "price": 75000,
    "image_url": "https://example.com/image.jpg",
    "category": "kaos",
    "stock": 50,
    "is_featured": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Product created successfully"
}
```

### 4. Update Product

**PUT** `/products/:id`

Update an existing product.

**Request Body:**

```json
{
  "name": "KAOS SABLON PREMIUM UPDATED",
  "description": "Kaos berkualitas tinggi dengan sablon premium yang sudah diupdate",
  "price": 80000,
  "image_url": "https://example.com/new-image.jpg",
  "category": "kaos",
  "stock": 45,
  "is_featured": true,
  "is_active": true
}
```

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "KAOS SABLON PREMIUM UPDATED",
    "description": "Kaos berkualitas tinggi dengan sablon premium yang sudah diupdate",
    "price": 80000,
    "image_url": "https://example.com/new-image.jpg",
    "category": "kaos",
    "stock": 45,
    "is_featured": true,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "message": "Product updated successfully"
}
```

### 5. Delete Product

**DELETE** `/products/:id`

Delete a product permanently.

**Response:**

```json
{
  "status": "success",
  "message": "Product deleted successfully"
}
```

### 6. Toggle Featured Status

**PATCH** `/products/:id/featured`

Toggle the featured status of a product.

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "KAOS SABLON PREMIUM",
    "description": "Kaos berkualitas tinggi dengan sablon premium",
    "price": 75000,
    "image_url": "https://example.com/image.jpg",
    "category": "kaos",
    "stock": 50,
    "is_featured": true,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "message": "Featured status updated successfully"
}
```

## Error Responses

### 400 Bad Request

```json
{
  "status": "error",
  "message": "Invalid input data",
  "error": "Detailed error message"
}
```

### 404 Not Found

```json
{
  "status": "error",
  "message": "Product not found"
}
```

### 500 Internal Server Error

```json
{
  "status": "error",
  "message": "Failed to process request",
  "error": "Detailed error message"
}
```

## Data Models

### Product

```json
{
  "id": "string (ObjectID)",
  "name": "string (required)",
  "description": "string",
  "price": "number (required, > 0)",
  "image_url": "string (URL)",
  "category": "string",
  "stock": "number (default: 0)",
  "is_featured": "boolean (default: false)",
  "is_active": "boolean (default: true)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Categories

Available product categories:

- `kaos` - Kaos/T-Shirt
- `jaket` - Jaket/Jacket
- `sweater` - Sweater
- `sport` - Sport Wear

## Usage Examples

### Frontend Integration

#### Load Featured Products

```javascript
const response = await fetch(
  "https://your-backend-url.com/products?featured=true"
);
const data = await response.json();
const featuredProducts = data.data;
```

#### Load All Products (Admin Dashboard)

```javascript
const response = await fetch("https://your-backend-url.com/products?all=true");
const data = await response.json();
const allProducts = data.data; // Includes active and inactive products
```

#### Load Active Products Only

```javascript
const response = await fetch("https://your-backend-url.com/products");
const data = await response.json();
const activeProducts = data.data; // Only active products
```

#### Create New Product

```javascript
const newProduct = {
  name: "KAOS SABLON PREMIUM",
  description: "Kaos berkualitas tinggi",
  price: 75000,
  category: "kaos",
  stock: 50,
  is_featured: true,
  is_active: true,
};

const response = await fetch("https://your-backend-url.com/products", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(newProduct),
});
```

#### Update Product

```javascript
const updateData = {
  name: "KAOS SABLON PREMIUM UPDATED",
  price: 80000,
  stock: 45,
};

const response = await fetch(
  `https://your-backend-url.com/products/${productId}`,
  {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updateData),
  }
);
```

#### Delete Product

```javascript
const response = await fetch(
  `https://your-backend-url.com/products/${productId}`,
  {
    method: "DELETE",
  }
);
```

#### Toggle Featured Status

```javascript
const response = await fetch(
  `https://your-backend-url.com/products/${productId}/featured`,
  {
    method: "PATCH",
  }
);
```
