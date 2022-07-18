//nolint:typecheck
package routes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRoom creates a room in the database
func CreateRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	if err := validate.Struct(room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to validate: %v", err)
		return
	}

	err := db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"uri": "/rooms/" + fmt.Sprint(room.ID)})
}

// GetRoom returns a room from the database
func GetRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
	}

	err = db.WithContext(ctx).First(&room, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// UpdateRoom updates a room in the database
func UpdateRoom(c *gin.Context) {
	var roomUpdate models.RoomUpdate
	var patchedRoom models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}
	if err := validate.Struct(roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to validate: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedRoom, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			log.Printf("Failed to find document: %v", err)
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not updated"})
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	err = db.WithContext(ctx).Model(&patchedRoom).Updates(roomUpdate.ToRoom()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not modified"})
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": patchedRoom})
}

// DeleteRoom deletes a room in the database
func DeleteRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		log.Printf("Failed to convert room ID: %v", err)
	}

	result := db.WithContext(ctx).Delete(&models.Room{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not deleted"})
		log.Printf("Failed to delete document: %v", result.Error)
		return
	} else if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		log.Printf("Failed to find document: %v", result.Error)
		return
	}

	c.JSON(http.StatusOK, nil)
}
