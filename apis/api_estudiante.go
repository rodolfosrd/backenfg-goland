package apis

import (
	"net/http"

	"github.com/202lp2/go2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//CRUD for items table
func EstudianteGET(c *gin.Context) {
	var lis []models.Estudiante
	db, _ := c.Get("db")
	conn := db.(gorm.DB)

	conn.Find(&lis)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Lista de los estudiantes",
		"r":   lis,
	})

}

func EstudiantePOST(c *gin.Context) {
	db, _ := c.Get("db")

	conn := db.(gorm.DB)

	var d models.Estudiante
	//d := models.Person{Name: c.PostForm("name"), Age: c.PostForm("age")}
	if err := c.BindJSON(&d); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//persona := models.Person{}
	//personas, _ := persona

	conn.Create(&d)
	c.JSON(http.StatusOK, &d)
}

func EstudianteGETID(c *gin.Context) {

	db, _ := c.Get("db")
	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.Estudiante
	if err := conn.First(&d, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &d)
}

func EstudianteUpdate(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.Estudiante
	if err := conn.First(&d, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	d.CODIGO = c.PostForm("nombre")
	c.BindJSON(&d)
	conn.Save(&d)
	c.JSON(http.StatusOK, &d)
}

func EstudianteDelete(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.Estudiante

	if err := conn.Where("id = ?", id).First(&d).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	conn.Unscoped().Delete(&d)
}
