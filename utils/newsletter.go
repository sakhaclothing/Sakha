package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendWelcomeEmail sends welcome email to new subscriber
func SendWelcomeEmail(email string) {
	subject := "Welcome to Sakha Clothing Newsletter!"
	content := fmt.Sprintf(`
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
	`)

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
	content := fmt.Sprintf(`
		<h2>New Product Alert! 🎉</h2>
		<p>We're excited to announce our latest product:</p>
		
		<div style="border: 1px solid #ddd; padding: 20px; margin: 20px 0; border-radius: 8px;">
			<h3>%s</h3>
			<p><strong>Price:</strong> Rp %,.0f</p>
			<p><strong>Category:</strong> %s</p>
			<p>%s</p>
		</div>
		
		<p>Be the first to get your hands on this amazing product!</p>
		<p><a href="https://sakhaclothing.shop/featuredproducts" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">View Product</a></p>
		
		<p>Best regards,<br>The Sakha Clothing Team</p>
	`, product.Name, product.Price, product.Category, product.Description)

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