package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendWelcomeEmail sends welcome email to new subscriber
func SendWelcomeEmail(email string) {
	subject := "Welcome to Sakha Clothing Newsletter!"
	content := `
		<h2>Welcome to Sakha Clothing!</h2>
		<p>Thank you for subscribing to our newsletter. You'll be the first to know about:</p>
		<ul>
			<li>New product releases</li>
			<li>Special offers and discounts</li>
			<li>Latest fashion trends</li>
			<li>Exclusive content</li>
		</ul>
		<p>Stay tuned for exciting updates!</p>
		<p>Best regards,<br>The Sakha Clothing Team</p>
	`

	err := SendEmail(email, subject, content)
	if err != nil {
		fmt.Printf("Failed to send welcome email to %s: %v\n", email, err)
	}
}

// SendNewProductNotificationToSubscribers sends notification to all active subscribers
func SendNewProductNotificationToSubscribers(product model.Product) {
	// Get all active subscribers
	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		fmt.Printf("Failed to fetch subscribers: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var subscribers []model.NewsletterSubscription
	if err := cursor.All(ctx, &subscribers); err != nil {
		fmt.Printf("Failed to decode subscribers: %v\n", err)
		return
	}

	// Create notification record
	notification := model.NewProductNotification{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Product:   product,
		SentAt:    time.Now(),
		SentCount: len(subscribers),
		Status:    "pending",
	}

	// Save notification record
	notificationCollection := config.DB.Collection("new_product_notifications")
	_, err = notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		fmt.Printf("Failed to save notification record: %v\n", err)
	}

	// Send email to each subscriber
	subject := fmt.Sprintf("New Product Alert: %s", product.Name)

	// Create product image HTML
	productImageHTML := ""
	if product.ImageURL != "" {
		productImageHTML = fmt.Sprintf(`
			<div style="text-align: center; margin: 20px 0;">
				<img src="%s" alt="%s" style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
			</div>
		`, product.ImageURL, product.Name)
	}

	content := fmt.Sprintf(`
		<h2>New Product Alert! 🎉</h2>
		<p>We're excited to announce our latest product:</p>
		
		<div style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9;">
			<h3 style="color: #333; margin-top: 0;">%s</h3>
			%s
			<p><strong>Price:</strong> <span style="color: #e74c3c; font-size: 18px;">%s</span></p>
			<p><strong>Category:</strong> <span style="background-color: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;">%s</span></p>
			<p style="color: #666; line-height: 1.6;">%s</p>
		</div>
		
		<p>Be the first to get your hands on this amazing product!</p>
		<p><a href="https://sakhaclothing.shop/featuredproducts" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">View Product</a></p>
		
		<p>Best regards,<br>The Sakha Clothing Team</p>
	`, product.Name, productImageHTML, formatRupiah(product.Price), product.Category, product.Description)

	successCount := 0
	for _, subscriber := range subscribers {
		err := SendEmail(subscriber.Email, subject, content)
		if err != nil {
			fmt.Printf("Failed to send notification to %s: %v\n", subscriber.Email, err)
		} else {
			successCount++
		}
	}

	// Update notification status
	_, err = notificationCollection.UpdateOne(
		ctx,
		bson.M{"_id": notification.ID},
		bson.M{
			"$set": bson.M{
				"status": "sent",
			},
		},
	)
	if err != nil {
		fmt.Printf("Failed to update notification status: %v\n", err)
	}

	fmt.Printf("Sent new product notification to %d/%d subscribers\n", successCount, len(subscribers))
}

// SendProductDeletedNotificationToSubscribers sends notification when a product is deleted
func SendProductDeletedNotificationToSubscribers(product model.Product) {
	// Get all active subscribers
	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		fmt.Printf("Failed to fetch subscribers: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var subscribers []model.NewsletterSubscription
	if err := cursor.All(ctx, &subscribers); err != nil {
		fmt.Printf("Failed to decode subscribers: %v\n", err)
		return
	}

	// Create notification record
	notification := model.NewProductNotification{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Product:   product,
		SentAt:    time.Now(),
		SentCount: len(subscribers),
		Status:    "pending",
	}

	// Save notification record
	notificationCollection := config.DB.Collection("product_deleted_notifications")
	_, err = notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		fmt.Printf("Failed to save deletion notification record: %v\n", err)
	}

	// Send email to each subscriber
	subject := fmt.Sprintf("Product Removed: %s", product.Name)

	// Create product image HTML
	productImageHTML := ""
	if product.ImageURL != "" {
		productImageHTML = fmt.Sprintf(`
			<div style="text-align: center; margin: 20px 0;">
				<img src="%s" alt="%s" style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); opacity: 0.7;">
			</div>
		`, product.ImageURL, product.Name)
	}

	content := fmt.Sprintf(`
		<h2>Product Removed Notice 📢</h2>
		<p>We want to inform you that the following product has been removed from our collection:</p>
		
		<div style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9; opacity: 0.8;">
			<h3 style="color: #333; margin-top: 0;">%s</h3>
			%s
			<p><strong>Price:</strong> <span style="color: #e74c3c; font-size: 18px;">%s</span></p>
			<p><strong>Category:</strong> <span style="background-color: #95a5a6; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;">%s</span></p>
			<p style="color: #666; line-height: 1.6;">%s</p>
		</div>
		
		<p>Don't worry! We have many other amazing products available. Check out our current collection!</p>
		<p><a href="https://sakhaclothing.shop/featuredproducts" style="background-color: #27ae60; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">Browse Products</a></p>
		
		<p>Best regards,<br>The Sakha Clothing Team</p>
	`, product.Name, productImageHTML, formatRupiah(product.Price), product.Category, product.Description)

	successCount := 0
	for _, subscriber := range subscribers {
		err := SendEmail(subscriber.Email, subject, content)
		if err != nil {
			fmt.Printf("Failed to send deletion notification to %s: %v\n", subscriber.Email, err)
		} else {
			successCount++
		}
	}

	// Update notification status
	_, err = notificationCollection.UpdateOne(
		ctx,
		bson.M{"_id": notification.ID},
		bson.M{
			"$set": bson.M{
				"status": "sent",
			},
		},
	)
	if err != nil {
		fmt.Printf("Failed to update deletion notification status: %v\n", err)
	}

	fmt.Printf("Sent product deletion notification to %d/%d subscribers\n", successCount, len(subscribers))
}

// SendProductUpdatedNotificationToSubscribers sends notification when a product is updated with significant changes
func SendProductUpdatedNotificationToSubscribers(product model.Product, oldProduct model.Product) {
	// Get all active subscribers
	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		fmt.Printf("Failed to fetch subscribers: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var subscribers []model.NewsletterSubscription
	if err := cursor.All(ctx, &subscribers); err != nil {
		fmt.Printf("Failed to decode subscribers: %v\n", err)
		return
	}

	// Check if there are significant changes worth notifying
	hasSignificantChanges := false
	changes := []string{}

	if product.Price != oldProduct.Price {
		hasSignificantChanges = true
		changes = append(changes, "price")
	}
	if product.Name != oldProduct.Name {
		hasSignificantChanges = true
		changes = append(changes, "name")
	}
	if product.ImageURL != oldProduct.ImageURL {
		hasSignificantChanges = true
		changes = append(changes, "image")
	}

	// Only send notification if there are significant changes
	if !hasSignificantChanges {
		return
	}

	// Create notification record
	notification := model.NewProductNotification{
		ID:        primitive.NewObjectID(),
		ProductID: product.ID,
		Product:   product,
		SentAt:    time.Now(),
		SentCount: len(subscribers),
		Status:    "pending",
	}

	// Save notification record
	notificationCollection := config.DB.Collection("product_updated_notifications")
	_, err = notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		fmt.Printf("Failed to save update notification record: %v\n", err)
	}

	// Send email to each subscriber
	subject := fmt.Sprintf("Product Updated: %s", product.Name)

	// Create product image HTML
	productImageHTML := ""
	if product.ImageURL != "" {
		productImageHTML = fmt.Sprintf(`
			<div style="text-align: center; margin: 20px 0;">
				<img src="%s" alt="%s" style="max-width: 300px; height: auto; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
			</div>
		`, product.ImageURL, product.Name)
	}

	// Create changes summary
	changesText := ""
	if len(changes) > 0 {
		changesText = fmt.Sprintf(`
			<div style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; margin: 10px 0; border-radius: 5px;">
				<p style="margin: 0; color: #856404;"><strong>What's New:</strong> %s</p>
			</div>
		`, strings.Join(changes, ", "))
	}

	content := fmt.Sprintf(`
		<h2>Product Updated! 🔄</h2>
		<p>We've updated one of our products with exciting changes:</p>
		
		<div style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px; background-color: #f9f9f9;">
			<h3 style="color: #333; margin-top: 0;">%s</h3>
			%s
			%s
			<p><strong>Price:</strong> <span style="color: #e74c3c; font-size: 18px;">%s</span></p>
			<p><strong>Category:</strong> <span style="background-color: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px;">%s</span></p>
			<p style="color: #666; line-height: 1.6;">%s</p>
		</div>
		
		<p>Check out the updated product now!</p>
		<p><a href="https://sakhaclothing.shop/featuredproducts" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">View Product</a></p>
		
		<p>Best regards,<br>The Sakha Clothing Team</p>
	`, product.Name, productImageHTML, changesText, formatRupiah(product.Price), product.Category, product.Description)

	successCount := 0
	for _, subscriber := range subscribers {
		err := SendEmail(subscriber.Email, subject, content)
		if err != nil {
			fmt.Printf("Failed to send update notification to %s: %v\n", subscriber.Email, err)
		} else {
			successCount++
		}
	}

	// Update notification status
	_, err = notificationCollection.UpdateOne(
		ctx,
		bson.M{"_id": notification.ID},
		bson.M{
			"$set": bson.M{
				"status": "sent",
			},
		},
	)
	if err != nil {
		fmt.Printf("Failed to update notification status: %v\n", err)
	}

	fmt.Printf("Sent product update notification to %d/%d subscribers\n", successCount, len(subscribers))
}

// Helper untuk format harga Rupiah
func formatRupiah(amount float64) string {
	return fmt.Sprintf("Rp %s", commaSeparator(int64(amount)))
}

// Helper untuk menambahkan tanda koma pada angka
func commaSeparator(n int64) string {
	in := fmt.Sprintf("%d", n)
	out := ""
	for i, v := range in {
		if i != 0 && (len(in)-i)%3 == 0 {
			out += "."
		}
		out += string(v)
	}
	return out
}
