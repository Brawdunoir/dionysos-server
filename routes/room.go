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

// CreateRoom godoc
// @Summary      Creates a room
// @Tags         Rooms
// @Accept       json
// @Param room body models.Room true "Room object"
// @Produce      json
// @Success      201 {object} utils.URIResponse "Room created"
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms [post]
func CreateRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err := db.WithContext(ctx).Create(&room).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Room not created"))
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, utils.CreateURIResponse("/rooms/"+fmt.Sprint(room.ID)))
}

// GetRoom godoc
// @Summary      Gets a room
// @Tags         Rooms
// @Param 			 id path int true "Room ID"
// @Produce      json
// @Success      200 {object} models.Room
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "Room not found"
// @Router       /rooms/{id} [get]
func GetRoom(c *gin.Context) {
	var room models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
	}

	err = db.WithContext(ctx).First(&room, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Room not found"))
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom godoc
// @Summary      Updates a room
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param id path int true "Room ID"
// @Param room body models.Room true "Room object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "Room not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [patch]
func UpdateRoom(c *gin.Context) {
	var roomUpdate models.RoomUpdate
	var patchedRoom models.Room

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedRoom, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Room not found"))
			return
		} else {
			c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Room not updated"))
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	err = db.WithContext(ctx).Model(&patchedRoom).Updates(roomUpdate.ToRoom()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Room not modified"))
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteRoom godoc
// @Summary      Deletes a room
// @Tags         Rooms
// @Param id path int true "Room ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "Room not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /rooms/{id} [delete]
func DeleteRoom(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid room ID"))
		log.Printf("Failed to convert room ID: %v", err)
	}

	result := db.WithContext(ctx).Delete(&models.Room{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Room not deleted"))
		log.Printf("Failed to delete document: %v", result.Error)
		return
	} else if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Room not found"))
		log.Printf("Failed to find document: %v", result.Error)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
