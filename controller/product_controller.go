package controller

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sakhaclothing/config"
	"github.com/sakhaclothing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get all products
func GetProducts(c *fiber.Ctx) error {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if we want featured products only
	featured := c.Query("featured")
	// Check if we want all products (for admin dashboard)
	all := c.Query("all")
	var filter bson.M

	if featured == "true" {
		// For featured products page - only active and featured
		filter = bson.M{"is_featured": true, "is_active": true}
	} else if all == "true" {
		// For admin dashboard - show all products (active and inactive)
		filter = bson.M{}
	} else {
		// Default - only active products
		filter = bson.M{"is_active": true}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch products",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode products",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   products,
	})
}

// Get product by ID
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	var product model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
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

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   product,
	})
}

// Create new product
func CreateProduct(c *fiber.Ctx) error {
	var product model.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if strings.TrimSpace(product.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Product name is required",
		})
	}

	if product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Product price must be greater than 0",
		})
	}

	// Set default values
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if product.Stock == 0 {
		product.Stock = 0
	}

	// Set default boolean values if not provided
	if !product.IsActive {
		product.IsActive = true
	}
	if !product.IsFeatured {
		product.IsFeatured = false
	}

	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create product",
			"error":   err.Error(),
		})
	}

	// Send notification to newsletter subscribers if product is active and featured
	if product.IsActive && product.IsFeatured {
		go sendNewProductNotificationToSubscribers(product)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"data":    product,
		"message": "Product created successfully",
	})
}

// Update product
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	// Check if product exists
	var existingProduct model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingProduct)
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

	// Parse update data
	var updateData model.Product
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Build update fields
	update := bson.M{"updated_at": time.Now()}

	if updateData.Name != "" {
		update["name"] = updateData.Name
	}
	if updateData.Description != "" {
		update["description"] = updateData.Description
	}
	if updateData.Price > 0 {
		update["price"] = updateData.Price
	}
	if updateData.ImageURL != "" {
		update["image_url"] = updateData.ImageURL
	}
	if updateData.Category != "" {
		update["category"] = updateData.Category
	}
	if updateData.Stock >= 0 {
		update["stock"] = updateData.Stock
	}

	// Update boolean fields
	update["is_featured"] = updateData.IsFeatured
	update["is_active"] = updateData.IsActive

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update product",
			"error":   err.Error(),
		})
	}

	// Fetch updated product
	var updatedProduct model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch updated product",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"data":    updatedProduct,
		"message": "Product updated successfully",
	})
}

// Delete product
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	// Check if product exists
	var product model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
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

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete product",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}

// Toggle featured status
func ToggleFeatured(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
		})
	}

	// Check if product exists
	var product model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
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

	// Toggle featured status
	newFeaturedStatus := !product.IsFeatured
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"is_featured": newFeaturedStatus,
			"updated_at":  time.Now(),
		}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update featured status",
			"error":   err.Error(),
		})
	}

	// Fetch updated product
	var updatedProduct model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch updated product",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"data":    updatedProduct,
		"message": "Featured status updated successfully",
	})
}
