package routes

import (
	"net/http"
	"os"

	utils "github.com/Brawdunoir/dionysos-server/utils/routes"
	"github.com/gin-gonic/gin"
)

// GetUser godoc
// @Summary      Get the api version
// @Tags         Misc
// @Produce      plain
// @Success      200 {string} plain "Version as x.x.x"
// @Failure      500 {object} utils.ErrorResponse "Version has not been set correctly in server side"
// @Router       /version [get]
func GetVersion(c *gin.Context) {
	var version string
	if version = os.Getenv("VERSION"); version == "" {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Version not set"))
		return
	}
	c.String(http.StatusOK, version)
}
