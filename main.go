package main

import (
	"fmt"
	"restApi-GoGin/config"
	"restApi-GoGin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db := config.LoadDatabase()
	config.RunMigration(db)

	router := gin.Default()
	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.AuthRouter(api)

	router.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
