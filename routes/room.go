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
	utils "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRoom creates a room in the database
func CreateRoom(c *gin.Context) {
	var owner models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&roomRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err := db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room not created"})
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, utils.CreateURIResponse("/rooms/"+fmt.Sprint(room.ID)))
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

	c.JSON(http.StatusOK, room)
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

	c.JSON(http.StatusNoContent, nil)
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

	c.JSON(http.StatusNoContent, nil)
}

// ConnectUserToRoom creates a room in the aganro database
func ConnectUserToRoom(c *gin.Context) {
	var userRequest ConnectUserRequest
	var user models.User
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	id := c.Param("id")

	colRooms, err := db.Collection(ctx, database.RoomsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}

	metaRoom, err := colRooms.ReadDocument(ctx, id, &room)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		log.Printf("Failed to get document: %v", err)
		return
	}

	colUsers, err := db.Collection(ctx, database.UsersCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot access database collection"})
		log.Printf("Failed to access collection: %v", err)
		return
	}
	metaUser, err := colUsers.ReadDocument(ctx, userRequest.UserID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		log.Printf("Failed to find document: %v", err)
		return
	}

	// add edge
	edgeCollection, _, err := graph.EdgeCollection(ctx, database.EdgeCollection)
	if err != nil {
		log.Fatalf("Failed to select edge collection: %v", err)
	}

	edge := Connection{From: string(metaUser.ID), To: string(metaRoom.ID)}
	fmt.Println(edge)
	_, err = edgeCollection.CreateDocument(ctx, edge)
	if err != nil {
		log.Fatalf("Failed to create edge document: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{"id": metaUser.ID})
}
