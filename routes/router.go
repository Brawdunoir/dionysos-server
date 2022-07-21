//nolint:typecheck
package routes

import (
	"net/http"

	"github.com/Brawdunoir/dionysos-server/constants"
	docs "github.com/Brawdunoir/dionysos-server/docs"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {
	basePath := constants.BasePath
	docs.SwaggerInfo.BasePath = basePath

	router.Use(options)

	r := router.Group(basePath)
	{
		// User
		userRouter := r.Group("/users")
		{
			userRouter.POST("", CreateUser)
			userRouter.GET("/:id", GetUser)
			userRouter.PATCH("/:id", UpdateUser)
			userRouter.DELETE("/:id", DeleteUser)
		}
		// Room
		roomRouter := r.Group("/rooms")
		{
			roomRouter.POST("", CreateRoom)
			roomRouter.GET("/:id", GetRoom)
			roomRouter.PATCH("/:id", UpdateRoom)
			roomRouter.DELETE("/:id", DeleteRoom)
		}
		// Version
		r.GET("/version", GetVersion)

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return router
}

// Middleware for CORS requests.
func options(c *gin.Context) {
	c.Header("Allow", "GET,POST,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "X-Requested-With, Authorization, Origin, Content-Type, Accept")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
