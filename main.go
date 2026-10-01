package main

import (
	"log"

	"example.com/rest_api/database"
	"example.com/rest_api/routes"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	err := utils.LoadJWTSecret()
	if err != nil {
		log.Fatal(err)
	}
	database.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8090") //localhost:8090
}
