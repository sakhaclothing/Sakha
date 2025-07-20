# Email Notifications with Product Images - Setup Guide

## Overview

Fitur email notifikasi produk telah ditingkatkan untuk menyertakan gambar produk dalam setiap email yang dikirim. Fitur ini memberikan pengalaman visual yang lebih baik kepada subscriber newsletter.

## What's New

### 🆕 **Enhanced Features:**

1. **Product Images in Emails**

   - Gambar produk ditampilkan dengan styling yang menarik
   - Responsive design untuk berbagai email client
   - Fallback handling jika gambar tidak tersedia

2. **Multiple Notification Types**

   - New Product Alerts
   - Product Update Notifications
   - Product Deletion Notifications
   - Featured Product Toggle Notifications

3. **Improved Email Design**
   - Color-coded notifications
   - Better typography and spacing
   - Professional call-to-action buttons
   - Mobile-responsive layout

## Setup Instructions

### 1. Backend Setup

#### Update Dependencies

Pastikan semua file yang diupdate sudah di-deploy:

```bash
# Files yang diupdate:
# - Sakha/utils/newsletter.go
# - Sakha/controller/product_controller.go
# - Sakha/model/product.go (sudah ada field ImageURL)
```

#### Database Collections

Sistem akan membuat collection baru secara otomatis:

- `product_deleted_notifications`
- `product_updated_notifications`

### 2. Email Configuration

#### SMTP Setup

Pastikan SMTP configuration sudah benar di database:

```javascript
// Collection: configurations
// Document ID: smtp
{
  "_id": "smtp",
  "SMTP_HOST": "smtp.gmail.com",
  "SMTP_PORT": "587",
  "SMTP_USERNAME": "your-email@gmail.com",
  "SMTP_PASSWORD": "your-app-password",
  "FROM_EMAIL": "your-email@gmail.com",
  "FROM_NAME": "Sakha Clothing"
}
```

#### Image URL Requirements

- URL harus accessible dari email client
- Format yang didukung: JPG, PNG, GIF, WebP
- Recommended size: 300px width, auto height
- Max file size: 5MB

### 3. Testing

#### Run Test Script

```bash
cd Sakha/scripts
go run test_email_with_images.go
```

#### Manual Testing

1. **Create New Product**

   ```bash
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

2. **Update Product**

   ```bash
   curl -X PUT https://your-backend-url.com/products/PRODUCT_ID \
     -H "Content-Type: application/json" \
     -d '{
       "price": 85000,
       "image_url": "https://example.com/new-image.jpg"
     }'
   ```

3. **Delete Product**

   ```bash
   curl -X DELETE https://your-backend-url.com/products/PRODUCT_ID
   ```

4. **Toggle Featured**
   ```bash
   curl -X PATCH https://your-backend-url.com/products/PRODUCT_ID/featured
   ```

## Email Templates

### 1. New Product Alert

- **Trigger**: Produk baru dengan `is_active: true` dan `is_featured: true`
- **Style**: Blue theme dengan gambar produk
- **Content**: Nama, harga, kategori, deskripsi, gambar

### 2. Product Update Notification

- **Trigger**: Perubahan signifikan (harga, nama, gambar)
- **Style**: Yellow highlight untuk perubahan
- **Content**: Gambar produk, perubahan yang dilakukan

### 3. Product Deletion Notification

- **Trigger**: Produk dihapus dari sistem
- **Style**: Gray theme dengan opacity
- **Content**: Gambar produk (dengan opacity), informasi produk

### 4. Featured Product Notification

- **Trigger**: Status featured diubah menjadi true
- **Style**: Blue theme
- **Content**: Gambar produk, informasi lengkap

## Configuration Options

### 1. Image Styling

```go
// Default image styling
style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);"
```

### 2. Color Themes

- **New Products**: `#007bff` (Blue)
- **Updates**: `#fff3cd` (Yellow)
- **Deletions**: `#95a5a6` (Gray)
- **Featured**: `#27ae60` (Green)

### 3. Notification Triggers

```go
// Significant changes that trigger notifications
- Price changes
- Name changes
- Image URL changes
- Featured status changes
- Product deletion
```

## Troubleshooting

### Issue: Images not displaying

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

### Issue: Duplicate notifications

**Solutions:**

1. Check notification logic
2. Verify trigger conditions
3. Review database records
4. Check async processing

## Best Practices

### 1. Image Optimization

- Use compressed images (JPEG for photos, PNG for graphics)
- Keep file size under 5MB
- Use CDN for better loading speed
- Provide alt text for accessibility

### 2. Email Content

- Keep subject lines clear and engaging
- Use emojis sparingly but effectively
- Include clear call-to-action buttons
- Test with multiple email clients

### 3. Performance

- Use async processing for notifications
- Implement rate limiting if needed
- Monitor email delivery rates
- Track notification success rates

## Monitoring

### 1. Database Collections

Monitor these collections for notification tracking:

- `new_product_notifications`
- `product_updated_notifications`
- `product_deleted_notifications`
- `newsletter_subscriptions`

### 2. Email Metrics

Track these metrics:

- Email delivery rate
- Open rate
- Click-through rate
- Bounce rate

### 3. Error Logging

Monitor these logs:

- SMTP errors
- Database errors
- Notification processing errors
- Image loading errors

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

## Support

### Documentation

- `Sakha/docs/product_email_notifications.md` - Detailed feature documentation
- `Sakha/docs/newsletter_feature.md` - Newsletter system overview
- `Sakha/docs/email_setup.md` - Email configuration guide

### Code Files

- `Sakha/utils/newsletter.go` - Email notification functions
- `Sakha/controller/product_controller.go` - Product CRUD with notifications
- `Sakha/scripts/test_email_with_images.go` - Test script

### Testing

- Run test script: `go run Sakha/scripts/test_email_with_images.go`
- Manual API testing with curl commands
- Email client compatibility testing
