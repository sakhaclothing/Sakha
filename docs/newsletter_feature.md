# Newsletter Feature - Sakha Clothing

## Overview

Fitur newsletter memungkinkan user untuk berlangganan update produk baru dan admin untuk mengirim notifikasi email kepada subscriber. Sistem ini terintegrasi dengan email service yang sudah ada.

## Features

### 🎯 **User Features:**

- ✅ Subscribe newsletter via homepage
- ✅ Unsubscribe via email link
- ✅ Welcome email otomatis
- ✅ Email validation
- ✅ Duplicate email handling

### 🎯 **Admin Features:**

- ✅ Dashboard newsletter management
- ✅ View all subscribers (active/inactive)
- ✅ Send manual notifications
- ✅ Auto-notification saat product baru
- ✅ Export subscribers to CSV
- ✅ Notification history tracking

### 🎯 **Technical Features:**

- ✅ Email template system
- ✅ Async email sending
- ✅ Database tracking
- ✅ API endpoints
- ✅ Frontend integration

## Database Collections

### 1. `newsletter_subscriptions`

```json
{
  "_id": "ObjectId",
  "email": "user@example.com",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 2. `new_product_notifications`

```json
{
  "_id": "ObjectId",
  "product_id": "ObjectId",
  "product": {
    "name": "Product Name",
    "price": 75000,
    "description": "Product description",
    "category": "kaos"
  },
  "sent_at": "2024-01-01T00:00:00Z",
  "sent_count": 150,
  "status": "sent"
}
```

## API Endpoints

### 1. Subscribe to Newsletter

**POST** `/newsletter/subscribe`

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Successfully subscribed to newsletter",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. Unsubscribe from Newsletter

**GET** `/newsletter/unsubscribe?email=user@example.com`

**Response:**

```json
{
  "status": "success",
  "message": "Successfully unsubscribed from newsletter"
}
```

### 3. Get All Subscribers (Admin)

**GET** `/newsletter/subscribers`

**Query Parameters:**

- `active` (optional): Filter by status

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "email": "user@example.com",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 4. Send New Product Notification

**POST** `/newsletter/notify/:productId`

**Response:**

```json
{
  "status": "success",
  "message": "New product notification sent to subscribers",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Product Name",
    "price": 75000,
    "description": "Product description"
  }
}
```

### 5. Get Notification History (Admin)

**GET** `/newsletter/history`

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "product_id": "507f1f77bcf86cd799439012",
      "product": {
        "name": "Product Name",
        "price": 75000,
        "description": "Product description"
      },
      "sent_at": "2024-01-01T00:00:00Z",
      "sent_count": 150,
      "status": "sent"
    }
  ]
}
```

## Frontend Integration

### 1. Homepage Newsletter Form

**File:** `sakhaclothing.github.io/index.html`

```html
<form class="newsletter__form" id="newsletterForm">
  <input
    type="email"
    placeholder="Enter your email"
    class="newsletter__input"
    id="newsletterEmail"
    required
  />
  <button type="submit" class="newsletter__button" id="newsletterButton">
    Subscribe Our product →
  </button>
</form>
<div
  id="newsletterMessage"
  class="newsletter__message"
  style="display: none;"
></div>
```

**JavaScript:** `sakhaclothing.github.io/script.js`

- Email validation
- API integration
- Success/error handling
- Loading states

### 2. Admin Dashboard

**File:** `dashboard/newsletter-management.html`

**Features:**

- Subscriber management
- Notification history
- Manual notification sending
- Export functionality

**JavaScript:** `dashboard/newsletter-management.js`

- Tab management
- Data loading
- CRUD operations
- CSV export

## Email Templates

### 1. Welcome Email

```html
<h2>Welcome to Sakha Clothing!</h2>
<p>
  Thank you for subscribing to our newsletter. You'll be the first to know
  about:
</p>
<ul>
  <li>New product releases</li>
  <li>Special offers and discounts</li>
  <li>Latest fashion trends</li>
  <li>Exclusive content</li>
</ul>
<p>Stay tuned for exciting updates!</p>
<p>Best regards,<br />The Sakha Clothing Team</p>
```

### 2. New Product Notification

```html
<h2>New Product Alert! 🎉</h2>
<p>We're excited to announce our latest product:</p>

<div
  style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px;"
>
  <h3>Product Name</h3>
  <p><strong>Price:</strong> Rp 75,000</p>
  <p><strong>Category:</strong> kaos</p>
  <p>Product description</p>
</div>

<p>Be the first to get your hands on this amazing product!</p>
<p><a href="https://your-website.com/featuredproducts">View Product</a></p>

<p>Best regards,<br />The Sakha Clothing Team</p>
```

## Auto-Notification System

### Product Creation Trigger

**File:** `Sakha/controller/product_controller.go`

```go
// Send notification to newsletter subscribers if product is active and featured
if product.IsActive && product.IsFeatured {
    go sendNewProductNotificationToSubscribers(product)
}
```

### Notification Process

1. **Product Created** → Check if active & featured
2. **Get Active Subscribers** → Query database
3. **Send Emails** → Async goroutine
4. **Track Results** → Save to database
5. **Update Status** → Mark as sent

## Setup Instructions

### 1. Backend Setup

```bash
# Ensure email configuration is set up
# Check Sakha/config/email_config.go

# Test email service
go run Sakha/scripts/test_email.go
```

### 2. Frontend Setup

```bash
# No additional setup required
# Newsletter form is already integrated in homepage
# Admin dashboard is accessible via sidebar
```

### 3. Database Setup

```bash
# Collections will be created automatically
# No manual setup required
```

## Testing

### 1. Test Newsletter Subscription

1. Buka homepage
2. Masukkan email di form newsletter
3. Cek email welcome
4. Cek database collection

### 2. Test Admin Dashboard

1. Login sebagai admin
2. Buka Newsletter Management
3. Cek subscriber list
4. Test send notification

### 3. Test Auto-Notification

1. Buat product baru dengan `is_active: true` dan `is_featured: true`
2. Cek email notification otomatis
3. Cek notification history

### 4. Test API Endpoints

```bash
# Subscribe
curl -X POST https://your-backend-url.com/newsletter/subscribe \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Get subscribers
curl https://your-backend-url.com/newsletter/subscribers

# Send notification
curl -X POST https://your-backend-url.com/newsletter/notify/productId
```

## Usage Examples

### 1. User Subscription Flow

```javascript
// User enters email in homepage
const email = "user@example.com";

// Frontend sends request
const response = await fetch("/newsletter/subscribe", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ email }),
});

// User receives welcome email
// User gets notifications for new products
```

### 2. Admin Notification Flow

```javascript
// Admin creates new product
const product = {
  name: "New T-Shirt",
  price: 75000,
  is_active: true,
  is_featured: true,
};

// System automatically sends notification
// All active subscribers receive email
```

### 3. Manual Notification

```javascript
// Admin selects product in dashboard
const productId = "507f1f77bcf86cd799439011";

// Admin clicks "Send Notification"
const response = await fetch(`/newsletter/notify/${productId}`, {
  method: "POST",
});

// All active subscribers receive email
```

## Troubleshooting

### Issue: Emails not sending

**Check:**

1. Email configuration in `config/email_config.go`
2. SMTP credentials
3. Network connectivity
4. Email service logs

### Issue: Subscribers not receiving notifications

**Check:**

1. Subscriber status (`is_active: true`)
2. Email template format
3. Email service status
4. Notification history

### Issue: Newsletter form not working

**Check:**

1. JavaScript console errors
2. API endpoint availability
3. CORS configuration
4. Network connectivity

### Issue: Admin dashboard not loading

**Check:**

1. Admin authentication
2. API endpoints
3. Database connectivity
4. JavaScript errors

## Security Considerations

### 1. Email Validation

- ✅ Frontend validation
- ✅ Backend validation
- ✅ Email format checking

### 2. Rate Limiting

- ✅ Prevent spam subscriptions
- ✅ Limit notification sending
- ✅ Protect API endpoints

### 3. Data Protection

- ✅ Email encryption
- ✅ Secure storage
- ✅ GDPR compliance

## Performance Optimization

### 1. Async Email Sending

- ✅ Goroutines for email sending
- ✅ Non-blocking operations
- ✅ Background processing

### 2. Database Optimization

- ✅ Indexed email field
- ✅ Efficient queries
- ✅ Connection pooling

### 3. Caching

- ✅ Subscriber count caching
- ✅ Template caching
- ✅ API response caching

## Future Enhancements

### 1. Advanced Features

- 📧 Email templates customization
- 📊 Analytics dashboard
- 🎯 Targeted campaigns
- 📱 SMS notifications

### 2. Integration

- 🔗 Social media integration
- 📈 Analytics integration
- 🛒 E-commerce integration
- 📊 CRM integration

### 3. Automation

- 🤖 AI-powered content
- ⏰ Scheduled campaigns
- 📅 Event-based triggers
- 🎨 Dynamic templates

## Related Files

### Backend

- `Sakha/model/newsletter.go` - Data models
- `Sakha/controller/newsletter_controller.go` - API handlers
- `Sakha/route/route.go` - Route definitions
- `Sakha/utils/email.go` - Email utilities

### Frontend

- `sakhaclothing.github.io/index.html` - Newsletter form
- `sakhaclothing.github.io/script.js` - Newsletter JavaScript
- `dashboard/newsletter-management.html` - Admin dashboard
- `dashboard/newsletter-management.js` - Admin JavaScript

### Documentation

- `Sakha/docs/newsletter_feature.md` - This file
- `Sakha/docs/email_setup.md` - Email setup guide
