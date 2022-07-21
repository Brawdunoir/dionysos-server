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

// CreateUser godoc
// @Summary      Creates a user
// @Tags         Users
// @Accept       json
// @Param user body models.User true "User object"
// @Produce      json
// @Success      201 {object} utils.URIResponse "User created"
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err := db.WithContext(ctx).Create(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("User not created"))
		log.Printf("Failed to create document: %v", err)
		return
	}

	c.JSON(http.StatusCreated, utils.CreateURIResponse("/users/"+fmt.Sprint(user.ID)))
}

// GetUser godoc
// @Summary      Gets a user
// @Tags         Users
// @Param 			 id path int true "User ID"
// @Produce      json
// @Success      200 {object} models.User
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	err = db.WithContext(ctx).First(&user, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("User not found"))
		log.Printf("Failed to find document: %v", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Updates a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param id path int true "User ID"
// @Param user body models.User true "User object"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	var userUpdate models.UserUpdate
	var patchedUser models.User

	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	// Test if data is valid
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse(err.Error()))
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	err = db.WithContext(ctx).First(&patchedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.CreateErrorResponse("User not found"))
			log.Printf("Failed to find document: %v", err)
			return
		} else {
			c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("User not updated"))
			log.Printf("Failed to modify document: %v", err)
			return
		}
	}

	err = db.WithContext(ctx).Model(&patchedUser).Updates(userUpdate.ToUser()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("User not modified"))
		log.Printf("Failed to modify document: %v", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteUser godoc
// @Summary      Deletes a user
// @Tags         Users
// @Param id path int true "User ID"
// @Success      204
// @Failure      400 {object} utils.ErrorResponse "Invalid request"
// @Failure      404 {object} utils.ErrorResponse "User not found"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, 1000*time.Millisecond)
	defer cancelCtx()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid user ID"))
		log.Printf("Failed to convert user ID: %v", err)
	}

	result := db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("User not deleted"))
		log.Printf("Failed to delete document: %v", result.Error)
		return
	} else if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("User not found"))
		log.Printf("Failed to find document: %v", result.Error)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
