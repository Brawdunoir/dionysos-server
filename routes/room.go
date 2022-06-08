//nolint:typecheck
package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Brawdunoir/dionysos-server/database"
	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// CreateRoom creates a room in the aganro database
func CreateRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		log.Printf("Failed to access collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room created", "id": meta.Key})
}

// GetRoom returns a room from the aganro database
func GetRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var result models.Room
	id := c.Param("id")

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		log.Printf("Failed to access collection: %v", err)
	}

	_, err = col.ReadDocument(ctx, id, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not found"})
		log.Printf("Failed to get document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room found", "room": result})
}

// UpdateRoom updates a room in the aganro database
func UpdateRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	var room models.Room
	id := c.Param("id")

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		log.Printf("Failed to access collection: %v", err)
	}

	patch := map[string]interface{}{
		"name": room.Name,
	}

	meta, err := col.UpdateDocument(ctx, id, patch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not modified"})
		log.Printf("Failed to update document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room modified", "id": meta.Key})
}

// DeleteRoom deletes a room in the aganro database
func DeleteRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()
	id := c.Param("id")

	col, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		log.Printf("Failed to access collection: %v", err)
	}

	meta, err := col.RemoveDocument(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Room not deleted"})
		log.Printf("Failed to delete document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted", "id": meta.Key})
}
