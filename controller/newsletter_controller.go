package controller

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"github.com/sakhaclothing/utils"
)

// Subscribe to newsletter
func SubscribeNewsletter(c *fiber.Ctx) error {
	var subscription model.NewsletterSubscription
	if err := c.BodyParser(&subscription); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate email
	if strings.TrimSpace(subscription.Email) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email is required",
		})
	}

	// Normalize email
	subscription.Email = strings.ToLower(strings.TrimSpace(subscription.Email))

	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email already exists
	var existingSubscription model.NewsletterSubscription
	err := collection.FindOne(ctx, bson.M{"email": subscription.Email}).Decode(&existingSubscription)
	if err == nil {
		// Email already exists, update to active
		_, err = collection.UpdateOne(
			ctx,
			bson.M{"email": subscription.Email},
			bson.M{
				"$set": bson.M{
					"is_active":  true,
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update subscription",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Email subscription reactivated successfully",
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to check existing subscription",
			"error":   err.Error(),
		})
	}

	// Create new subscription
	subscription.ID = primitive.NewObjectID()
	subscription.IsActive = true
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	_, err = collection.InsertOne(ctx, subscription)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create subscription",
			"error":   err.Error(),
		})
	}

	// Send welcome email
	go utils.SendWelcomeEmail(subscription.Email)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully subscribed to newsletter",
		"data":    subscription,
	})
}

// Unsubscribe from newsletter
func UnsubscribeNewsletter(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email parameter is required",
		})
	}

	email = strings.ToLower(strings.TrimSpace(email))

	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to unsubscribe",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Email not found in subscriptions",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully unsubscribed from newsletter",
	})
}

// Get all subscribers (admin only)
func GetSubscribers(c *fiber.Ctx) error {
	collection := config.DB.Collection("newsletter_subscriptions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get filter from query
	active := c.Query("active")
	var filter bson.M

	if active == "true" {
		filter = bson.M{"is_active": true}
	} else if active == "false" {
		filter = bson.M{"is_active": false}
	} else {
		filter = bson.M{}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch subscribers",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var subscribers []model.NewsletterSubscription
	if err := cursor.All(ctx, &subscribers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode subscribers",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   subscribers,
	})
}

// Send notification for new product
func SendNewProductNotification(c *fiber.Ctx) error {
	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Product ID is required",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	// Get product details
	productCollection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product model.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch product",
			"error":   err.Error(),
		})
	}

	// Send notification to all active subscribers
	go utils.SendNewProductNotificationToSubscribers(product)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "New product notification sent to subscribers",
		"data":    product,
	})
}

// Get notification history (admin only)
func GetNotificationHistory(c *fiber.Ctx) error {
	collection := config.DB.Collection("new_product_notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch notifications",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var notifications []model.NewProductNotification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode notifications",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   notifications,
	})
}
