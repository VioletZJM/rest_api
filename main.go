package main

import (
	"example.com/rest_api/database"
	"example.com/rest_api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8090") //localhost:8090
}
