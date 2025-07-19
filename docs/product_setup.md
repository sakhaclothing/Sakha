# Product Management Setup Guide

## Overview

Fitur Product Management memungkinkan admin untuk mengelola produk-produk yang ditampilkan di website Sakha Clothing. Admin dapat menambah, mengedit, menghapus, dan mengatur status featured produk.

## Features

- ✅ CRUD (Create, Read, Update, Delete) produk
- ✅ Toggle status featured produk
- ✅ Filter produk berdasarkan kategori dan status
- ✅ Search produk berdasarkan nama dan deskripsi
- ✅ Integrasi dengan featured products page
- ✅ Responsive admin dashboard

## File Structure

### Backend (Go + MongoDB)

```
Sakha/
├── model/
│   └── product.go              # Product data model
├── controller/
│   └── product_controller.go   # Product CRUD operations
├── route/
│   └── route.go               # API routes (updated)
└── docs/
    ├── product_api.md         # API documentation
    └── product_setup.md       # This file
```

### Frontend

```
dashboard/
├── product-management.html    # Admin product management page
└── product-management.js      # Product management JavaScript

featuredproducts/
├── index.html                # Featured products page (updated)
└── script.js                 # Dynamic product loading (updated)
```

## API Endpoints

| Method | Endpoint                  | Description                |
| ------ | ------------------------- | -------------------------- |
| GET    | `/products`               | Get all products           |
| GET    | `/products?featured=true` | Get featured products only |
| GET    | `/products/:id`           | Get product by ID          |
| POST   | `/products`               | Create new product         |
| PUT    | `/products/:id`           | Update product             |
| DELETE | `/products/:id`           | Delete product             |
| PATCH  | `/products/:id/featured`  | Toggle featured status     |

## Setup Instructions

### 1. Backend Setup

#### Update Dependencies

Pastikan MongoDB driver sudah terinstall:

```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/bson/primitive
```

#### Database Migration

Tidak perlu migration khusus karena MongoDB menggunakan schema-less design. Collection `products` akan dibuat otomatis saat pertama kali insert data.

#### Environment Variables

Pastikan environment variable `MONGOSTRING` sudah diset dengan benar di file `.env` atau environment system.

### 2. Frontend Setup

#### Update API Base URL

Update `apiBaseUrl` di file JavaScript:

- `dashboard/product-management.js` (line 4)
- `featuredproducts/script.js` (line 3)

Ganti `https://sakhaclothing.shop` dengan URL backend yang sebenarnya.

#### Access Admin Dashboard

1. Buka `/dashboard/product-management.html`
2. Login sebagai admin
3. Mulai mengelola produk

### 3. Testing

#### Test API Endpoints

```bash
# Get all products
curl -X GET https://your-backend-url.com/products

# Create new product
curl -X POST https://your-backend-url.com/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "KAOS SABLON PREMIUM",
    "description": "Kaos berkualitas tinggi",
    "price": 75000,
    "category": "kaos",
    "stock": 50,
    "is_featured": true,
    "is_active": true
  }'

# Get featured products
curl -X GET "https://your-backend-url.com/products?featured=true"
```

#### Test Frontend

1. Buka `/dashboard/product-management.html`
2. Coba tambah produk baru
3. Edit produk yang sudah ada
4. Toggle status featured
5. Hapus produk
6. Test search dan filter

## Usage Guide

### Admin Dashboard

#### Adding New Product

1. Klik tombol "Add Product"
2. Isi form dengan data produk:
   - **Name**: Nama produk (wajib)
   - **Description**: Deskripsi produk
   - **Price**: Harga dalam Rupiah (wajib)
   - **Category**: Pilih kategori
   - **Image URL**: URL gambar produk
   - **Stock**: Jumlah stok
   - **Active**: Status aktif/nonaktif
   - **Featured**: Status featured/non-featured
3. Klik "Save Product"

#### Editing Product

1. Klik icon edit (✏️) pada produk yang ingin diedit
2. Update data yang diperlukan
3. Klik "Save Product"

#### Toggle Featured Status

1. Klik tombol "Featured" atau "Not Featured" pada produk
2. Status akan berubah secara otomatis

#### Deleting Product

1. Klik icon delete (🗑️) pada produk
2. Konfirmasi penghapusan
3. Produk akan dihapus secara permanen

#### Search and Filter

- **Search**: Ketik nama atau deskripsi produk
- **Category Filter**: Filter berdasarkan kategori
- **Status Filter**: Filter berdasarkan status (Active/Inactive/Featured)

### Featured Products Page

Halaman `/featuredproducts/` akan otomatis menampilkan produk yang memiliki status `is_featured: true` dan `is_active: true`.

Jika API tidak tersedia, halaman akan menampilkan produk fallback yang sudah didefinisikan di `script.js`.

## Data Model

### Product Schema

```json
{
  "_id": "ObjectId",
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

### Categories

- `kaos` - Kaos/T-Shirt
- `jaket` - Jaket/Jacket
- `sweater` - Sweater
- `sport` - Sport Wear

## Troubleshooting

### Common Issues

#### 1. API Connection Error

**Problem**: Frontend tidak bisa connect ke backend
**Solution**:

- Periksa `apiBaseUrl` di file JavaScript
- Pastikan backend server berjalan
- Periksa CORS settings

#### 2. MongoDB Connection Error

**Problem**: Backend tidak bisa connect ke MongoDB
**Solution**:

- Periksa `MONGOSTRING` environment variable
- Pastikan MongoDB server berjalan
- Periksa network connectivity

#### 3. Product Not Showing in Featured Page

**Problem**: Produk tidak muncul di halaman featured
**Solution**:

- Pastikan `is_featured: true` dan `is_active: true`
- Periksa API response di browser developer tools
- Cek console untuk error messages

#### 4. Image Not Loading

**Problem**: Gambar produk tidak muncul
**Solution**:

- Periksa URL gambar di database
- Pastikan URL accessible dari internet
- Gunakan placeholder image jika URL tidak valid

### Debug Mode

Untuk debugging, buka browser developer tools dan periksa:

- **Console**: Error messages
- **Network**: API requests dan responses
- **Application**: Local storage dan session storage

## Security Considerations

1. **Authentication**: Pastikan admin dashboard dilindungi dengan authentication
2. **Authorization**: Hanya admin yang bisa mengakses product management
3. **Input Validation**: Backend sudah memvalidasi input data
4. **SQL Injection**: Tidak relevan untuk MongoDB, tapi tetap validasi input
5. **XSS Protection**: Sanitasi output HTML

## Performance Optimization

1. **Pagination**: Untuk produk yang banyak, implementasi pagination
2. **Caching**: Cache featured products untuk performa lebih baik
3. **Image Optimization**: Compress dan resize gambar produk
4. **CDN**: Gunakan CDN untuk gambar produk

## Future Enhancements

1. **Bulk Operations**: Import/export produk dalam batch
2. **Image Upload**: Upload gambar langsung dari admin dashboard
3. **Product Variants**: Support untuk ukuran dan warna
4. **Inventory Management**: Tracking stok otomatis
5. **Analytics**: Dashboard analytics untuk produk
6. **SEO**: Meta tags dan URL optimization untuk produk
