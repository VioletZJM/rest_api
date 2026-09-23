package main

import (
	"log"

	"example.com/rest_api/database"
	"example.com/rest_api/routes"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	err := utils.LoadJWTSecret()
	if err != nil {
		log.Fatal(err)
	}

	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8090") //localhost:8090
}
