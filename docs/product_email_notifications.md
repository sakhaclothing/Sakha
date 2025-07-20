# Product Email Notifications with Images - Sakha Clothing

## Overview

Sistem notifikasi email produk telah ditingkatkan untuk menyertakan gambar produk dalam setiap email yang dikirim. Fitur ini memastikan subscriber mendapatkan informasi visual yang lengkap tentang produk baru, update, atau yang dihapus.

## Features

### 🎯 **Enhanced Email Notifications:**

- ✅ **Product Images**: Setiap email menyertakan gambar produk
- ✅ **New Product Alerts**: Notifikasi otomatis saat produk baru ditambahkan
- ✅ **Product Updates**: Notifikasi saat produk diupdate dengan perubahan signifikan
- ✅ **Product Deletion**: Notifikasi saat produk dihapus
- ✅ **Featured Product Toggle**: Notifikasi saat produk menjadi featured
- ✅ **Responsive Email Design**: Email yang responsif dan menarik

### 🎯 **Notification Types:**

1. **New Product Notification**

   - Trigger: Produk baru dibuat dengan `is_active: true` dan `is_featured: true`
   - Content: Gambar produk, nama, harga, kategori, deskripsi
   - Style: Background hijau, tombol "View Product"

2. **Product Update Notification**

   - Trigger: Perubahan signifikan (harga, nama, gambar)
   - Content: Gambar produk, perubahan yang dilakukan, informasi produk
   - Style: Background kuning untuk highlight perubahan

3. **Product Deletion Notification**

   - Trigger: Produk dihapus dari sistem
   - Content: Gambar produk (dengan opacity), informasi produk yang dihapus
   - Style: Background abu-abu, tombol "Browse Products"

4. **Featured Product Notification**
   - Trigger: Status featured diubah menjadi true
   - Content: Gambar produk, informasi lengkap produk
   - Style: Background biru, tombol "View Product"

## Email Templates

### 1. New Product Alert

```html
<h2>New Product Alert! 🎉</h2>
<p>We're excited to announce our latest product:</p>

<div
  style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9;"
>
  <h3 style="color: #333; margin-top: 0;">Product Name</h3>

  <!-- Product Image -->
  <div style="text-align: center; margin: 20px 0;">
    <img
      src="PRODUCT_IMAGE_URL"
      alt="Product Name"
      style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);"
    />
  </div>

  <p>
    <strong>Price:</strong>
    <span style="color: #e74c3c; font-size: 18px;">Rp 75,000</span>
  </p>
  <p>
    <strong>Category:</strong>
    <span
      style="background-color: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;"
      >kaos</span
    >
  </p>
  <p style="color: #666; line-height: 1.6;">Product description</p>
</div>

<p>Be the first to get your hands on this amazing product!</p>
<p>
  <a
    href="https://sakhaclothing.shop/featuredproducts"
    style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;"
    >View Product</a
  >
</p>
```

### 2. Product Update Notification

```html
<h2>Product Updated! 🔄</h2>
<p>We've updated one of our products with exciting changes:</p>

<div
  style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9;"
>
  <h3 style="color: #333; margin-top: 0;">Product Name</h3>

  <!-- Product Image -->
  <div style="text-align: center; margin: 20px 0;">
    <img
      src="PRODUCT_IMAGE_URL"
      alt="Product Name"
      style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);"
    />
  </div>

  <!-- Changes Summary -->
  <div
    style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; margin: 10px 0; border-radius: 5px;"
  >
    <p style="margin: 0; color: #856404;">
      <strong>What's New:</strong> price, image
    </p>
  </div>

  <p>
    <strong>Price:</strong>
    <span style="color: #e74c3c; font-size: 18px;">Rp 75,000</span>
  </p>
  <p>
    <strong>Category:</strong>
    <span
      style="background-color: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;"
      >kaos</span
    >
  </p>
  <p style="color: #666; line-height: 1.6;">Product description</p>
</div>

<p>Check out the updated product now!</p>
<p>
  <a
    href="https://sakhaclothing.shop/featuredproducts"
    style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;"
    >View Product</a
  >
</p>
```

### 3. Product Deletion Notification

```html
<h2>Product Removed Notice 📢</h2>
<p>
  We want to inform you that the following product has been removed from our
  collection:
</p>

<div
  style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9; opacity: 0.8;"
>
  <h3 style="color: #333; margin-top: 0;">Product Name</h3>

  <!-- Product Image (with opacity) -->
  <div style="text-align: center; margin: 20px 0;">
    <img
      src="PRODUCT_IMAGE_URL"
      alt="Product Name"
      style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); opacity: 0.7;"
    />
  </div>

  <p>
    <strong>Price:</strong>
    <span style="color: #e74c3c; font-size: 18px;">Rp 75,000</span>
  </p>
  <p>
    <strong>Category:</strong>
    <span
      style="background-color: #95a5a6; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;"
      >kaos</span
    >
  </p>
  <p style="color: #666; line-height: 1.6;">Product description</p>
</div>

<p>
  Don't worry! We have many other amazing products available. Check out our
  current collection!
</p>
<p>
  <a
    href="https://sakhaclothing.shop/featuredproducts"
    style="background-color: #27ae60; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;"
    >Browse Products</a
  >
</p>
```

## Implementation Details

### 1. Image Handling

```go
// Create product image HTML
productImageHTML := ""
if product.ImageURL != "" {
    productImageHTML = fmt.Sprintf(`
        <div style="text-align: center; margin: 20px 0;">
            <img src="%s" alt="%s" style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
        </div>
    `, product.ImageURL, product.Name)
}
```

### 2. Notification Triggers

#### New Product Creation

```go
// In CreateProduct function
if product.IsActive && product.IsFeatured {
    go utils.SendNewProductNotificationToSubscribers(product)
}
```

#### Product Update

```go
// In UpdateProduct function
// Send notification if product becomes featured and active
if updatedProduct.IsFeatured && updatedProduct.IsActive && !existingProduct.IsFeatured {
    go utils.SendNewProductNotificationToSubscribers(updatedProduct)
}

// Send notification for significant product updates
go utils.SendProductUpdatedNotificationToSubscribers(updatedProduct, existingProduct)
```

#### Product Deletion

```go
// In DeleteProduct function
// Send notification to newsletter subscribers about product deletion
go utils.SendProductDeletedNotificationToSubscribers(product)
```

#### Featured Toggle

```go
// In ToggleFeatured function
// Send notification if product becomes featured and is active
if newFeaturedStatus && updatedProduct.IsActive {
    go utils.SendNewProductNotificationToSubscribers(updatedProduct)
}
```

### 3. Database Collections

#### New Collections Added:

- `product_deleted_notifications` - Track deletion notifications
- `product_updated_notifications` - Track update notifications

#### Existing Collections:

- `new_product_notifications` - Track new product notifications
- `newsletter_subscriptions` - Subscriber data

## Email Styling Features

### 1. Responsive Design

- Max-width: 300px untuk gambar
- Auto height untuk maintain aspect ratio
- Border radius dan shadow untuk estetika

### 2. Color Coding

- **New Products**: Blue theme (#007bff)
- **Updates**: Yellow highlight (#fff3cd)
- **Deletions**: Gray theme (#95a5a6)
- **Featured**: Green theme (#27ae60)

### 3. Visual Elements

- Product images dengan border radius dan shadow
- Price highlighting dengan warna merah
- Category badges dengan background color
- Call-to-action buttons yang menarik

## Testing

### 1. Test New Product Notification

```bash
# Create a new product with is_active: true and is_featured: true
curl -X POST https://your-backend-url.com/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "price": 75000,
    "description": "Test description",
    "image_url": "https://example.com/image.jpg",
    "category": "kaos",
    "is_active": true,
    "is_featured": true
  }'
```

### 2. Test Product Update Notification

```bash
# Update product with significant changes
curl -X PUT https://your-backend-url.com/products/PRODUCT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "price": 85000,
    "image_url": "https://example.com/new-image.jpg"
  }'
```

### 3. Test Product Deletion Notification

```bash
# Delete a product
curl -X DELETE https://your-backend-url.com/products/PRODUCT_ID
```

### 4. Test Featured Toggle Notification

```bash
# Toggle featured status
curl -X PATCH https://your-backend-url.com/products/PRODUCT_ID/featured
```

## Configuration

### 1. Image URL Requirements

- URL harus accessible dari email client
- Format yang didukung: JPG, PNG, GIF, WebP
- Recommended size: 300px width, auto height
- Max file size: 5MB

### 2. Email Client Compatibility

- Gmail: ✅ Fully supported
- Outlook: ✅ Supported
- Apple Mail: ✅ Supported
- Yahoo Mail: ✅ Supported
- Mobile clients: ✅ Responsive

### 3. Performance Considerations

- Images loaded asynchronously
- Fallback text for images
- Optimized HTML structure
- Minimal CSS for compatibility

## Troubleshooting

### Issue: Images not displaying in emails

**Solutions:**

1. Check image URL accessibility
2. Verify image format compatibility
3. Test with different email clients
4. Check image size limits

### Issue: Email layout broken

**Solutions:**

1. Test with inline CSS
2. Check email client compatibility
3. Verify HTML structure
4. Test responsive design

### Issue: Notifications not sending

**Solutions:**

1. Check SMTP configuration
2. Verify subscriber list
3. Check notification triggers
4. Review error logs

## Future Enhancements

### 1. Advanced Features

- 📧 Custom email templates
- 🎨 Dynamic color schemes
- 📱 Mobile-optimized layouts
- 🌐 Multi-language support

### 2. Analytics

- 📊 Email open rates
- 🔗 Click tracking
- 📈 Conversion metrics
- 📋 A/B testing

### 3. Personalization

- 👤 Personalized recommendations
- 🎯 Targeted campaigns
- 📅 Scheduled notifications
- 🛒 Purchase history integration

## Related Files

### Backend

- `Sakha/utils/newsletter.go` - Email notification functions
- `Sakha/controller/product_controller.go` - Product CRUD with notifications
- `Sakha/model/product.go` - Product data model
- `Sakha/utils/email.go` - Email sending utilities

### Documentation

- `Sakha/docs/newsletter_feature.md` - Newsletter system overview
- `Sakha/docs/email_setup.md` - Email configuration guide
- `Sakha/docs/product_api.md` - Product API documentation
