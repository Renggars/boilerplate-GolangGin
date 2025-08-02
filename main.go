package main

import (
	"fmt"
	"log"
	"restApi-GoGin/config"
	"restApi-GoGin/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	router.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
