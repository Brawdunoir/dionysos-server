package main

import (
	"os"

	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/gin-gonic/gin"
)

// VERSION as a constant variable
const VERSION = "0.1.0"

// @version         0.1.0

// @title           Dionysos
// @description     API instance for dionysos client application.

// @contact.name   API Support
// @contact.url    https://github.com/Brawdunoir/dionysos-server/issues

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	os.Setenv("VERSION", VERSION)

	router := routes.SetupRouter(gin.Default())

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
