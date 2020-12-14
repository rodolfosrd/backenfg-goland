package routers

import (
	"net/http"
	"strings"

	"github.com/202lp2/go2/apis"
	"github.com/202lp2/go2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {

	conn, err := connectDBmysql()
	if err != nil {
		panic("failed to connect database: " + err.Error())
		//return
	}
	// Migrate the schema
	conn.AutoMigrate(
		&models.Person{},
		&models.User{},
		//&models.Estudiante{},
	)

	r := gin.Default()

	//config := cors.DefaultConfig() https://github.com/rs/cors
	//config.AllowOrigins = []string{"http://localhost", "http://localhost:8086"}

	r.Use(CORSMiddleware())

	r.Use(dbMiddleware(*conn))

	v1 := r.Group("/v1")
	{
		v1.GET("/ping", apis.ItemsIndex)

		v1.GET("/estudiante", apis.EstudianteGET)
		v1.POST("/estudiante", authMiddleWare(), apis.EstudiantePOST)
		v1.GET("/estudiante/:id", apis.EstudianteGETID)
		v1.PUT("/estudiante/:id", apis.EstudianteUpdate)
		v1.DELETE("/estudiante/:id", apis.EstudianteDelete)

		v1.GET("/persons", apis.PersonsIndex)
		v1.POST("/persons", authMiddleWare(), apis.PersonsCreate)
		v1.GET("/persons/:id", apis.PersonsGet)
		v1.PUT("/persons/:id", apis.PersonsUpdate)
		v1.DELETE("/persons/:id", apis.PersonsDelete)

		v1.GET("/course", apis.CourseIndex)
		v1.POST("/course", authMiddleWare(), apis.CourseCreate)
		v1.GET("/course/:id", apis.CourseGet)
		v1.PUT("/course/:id", apis.CourseUpdate)
		v1.DELETE("/course/:id", apis.CourseDelete)

		v1.GET("/educationlevel", apis.EducationLevelIndex)
		v1.POST("/educationlevel", authMiddleWare(), apis.EducationLevelCreate)
		v1.GET("/educationlevel/:id", apis.EducationLevelGet)
		v1.PUT("/educationlevel/:id", apis.EducationLevelUpdate)
		v1.DELETE("/educationlevel/:id", apis.EducationLevelDelete)

		v1.GET("/academicperiod", apis.AcademicPeriodIndex)
		v1.POST("/academicperiod", authMiddleWare(), apis.AcademicPeriodCreate)
		v1.GET("/academicperiod/:id", apis.AcademicPeriodGet)
		v1.PUT("/academicperiod/:id", apis.AcademicPeriodUpdate)
		v1.DELETE("/academicperiod/:id", apis.AcademicPeriodDelete)

		v1.GET("/academicplan", apis.AcademicPlanIndex)
		v1.POST("/academicplan", authMiddleWare(), apis.AcademicPlanCreate)
		v1.GET("/academicplan/:id", apis.AcademicPlanGet)
		v1.PUT("/academicplan/:id", apis.AcademicPlanUpdate)
		v1.DELETE("/academicplan/:id", apis.AcademicPlanDelete)

		v1.GET("/users", apis.UsersIndex)
		v1.POST("/users", apis.UsersCreate)
		v1.GET("/users/:id", apis.UsersGet)
		v1.PUT("/users/:id", apis.UsersUpdate)
		v1.DELETE("/users/:id", apis.UsersDelete)
		v1.POST("/login", apis.UsersLogin)
		v1.POST("/logout", apis.UsersLogout)
	}
	return r
}

func connectDBmysql() (c *gorm.DB, err error) {

	//dsn := "docker:docker@tcp(mysql-db:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:Thiagosrd150@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open("mysql", "root:pass123@localhost/dbcontacts?charset=utf8&parseTime=True&loc=Local")
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Error al conectar a la Basde de Datos de Mysql: " + err.Error())
	}

	return conn, err
}

func connectDB() (c *gorm.DB, err error) {
	dsn := "user=jrnfnpgtnbqvnx password=50ff322d60f7d2089ccda0a687924c18098675b5f60d8150f24e72fa1929716c host=ec2-23-23-36-227.compute-1.amazonaws.com	dbname=d603guqgrv4va2 port=5432 sslmode=require TimeZone=Asia/Shanghai"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la Basde de Datos de Postgress: " + err.Error())
	}
	return conn, err
}

func dbMiddleware(conn gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//c.Header("Access-Control-Allow-Origin", "http://localhost, http://localhost:8086,")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE ")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//https://dev.to/stevensunflash/a-working-solution-to-jwt-creation-and-invalidation-in-golang-4oe4

//https://www.nexmo.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr
func authMiddleWare() gin.HandlerFunc { //ExtractToken
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]
		//fmt.Printf("Bearer (%v) \n", token)
		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated (IsTokenValid)."})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
