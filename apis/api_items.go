package apis

import (
	"net/http"

	"github.com/202lp2/go2/models"
	"github.com/gin-gonic/gin"
)

//CRUD for items table
func ItemsIndex(c *gin.Context) {
	s := models.Item{Title: "Sean", Notes: "nnn"}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hola we " + s.Title,
	})
}
