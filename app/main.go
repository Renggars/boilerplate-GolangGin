package main

import (
	"fmt"
	"log"

	_ "restApi-GoGin/docs"
	"restApi-GoGin/src/config"
	"restApi-GoGin/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Boilerplate Go Gin API
// @version         1.0
// @description     A boilerplate REST API using Go and Gin framework with authentication system

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	config.LoadConfig()
	db := config.LoadDatabase()
	config.RunMigration(db)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		// AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		// AllowHeaders:     []string{"Authorization", "Content-Type"},
		// MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.AuthRouter(api)
	routes.UserRouter(api)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	router.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
