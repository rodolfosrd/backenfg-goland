package apis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/202lp2/go2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//CRUD for items table
func UsersIndex(c *gin.Context) {
	var lis []models.User

	db, _ := c.Get("db")

	conn := db.(gorm.DB)
	// Migrate the schema
	//conn.AutoMigrate(&models.User{})

	conn.Find(&lis)
	c.JSON(http.StatusOK, gin.H{
		"msg": "thank you",
		"r":   lis,
	})

}

func UsersCreate(c *gin.Context) {
	db, _ := c.Get("db")

	conn := db.(gorm.DB)

	////var d models.User
	////d := models.User{Email: c.PostForm("email")} //, PasswordHash: c.PostForm("password_hash")
	//if err := c.BindJSON(&d); err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}
	//d.PasswordHash = string("123")
	//conn.Create(&d)
	//c.JSON(http.StatusOK, &d)
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Register(&conn)
	if err != nil {
		fmt.Println("Error in user.Register()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": user.ID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
	})
}

func UsersGet(c *gin.Context) {

	db, _ := c.Get("db")

	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.User
	if err := conn.First(&d, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &d)
}

func UsersUpdate(c *gin.Context) {
	db, _ := c.Get("db")

	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.User
	if err := conn.First(&d, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	d.Email = c.PostForm("email")
	d.PasswordHash = c.PostForm("password_hash")
	c.BindJSON(&d)
	conn.Save(&d)
	c.JSON(http.StatusOK, &d)
}

func UsersDelete(c *gin.Context) {
	db, _ := c.Get("db")

	conn := db.(gorm.DB)

	id := c.Param("id")
	var d models.User

	if err := conn.Where("id = ?", id).First(&d).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	conn.Unscoped().Delete(&d)
}

func UsersLogin(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(gorm.DB)

	err = user.IsAuthenticated(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": user.ID,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "There was an error authenticating.",
	})
}

func UsersLogout(c *gin.Context) { //ToDo

	bearer := c.Request.Header.Get("Authorization")
	split := strings.Split(bearer, "Bearer ")
	if len(split) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
		c.Abort()
		return
	}
	token := split[1]
	//fmt.Printf("Bearer (%v) \n", token)
	isValid := models.DelTokenValid(token)
	if isValid == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated (IsTokenValid)."})
		c.Abort()
	} else {
		c.Set("user_id", nil)
		c.Next()
	}
	/*
		au, err := ExtractTokenMetadata(c.Request)
		if err != nil {
		   c.JSON(http.StatusUnauthorized, "unauthorized")
		   return
		}
		deleted, delErr := DeleteAuth(au.AccessUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
		   c.JSON(http.StatusUnauthorized, "unauthorized")
		   return
		}
	*/
	c.JSON(http.StatusOK, "Successfully logged out ")
}
