package main

import (
	"github.com/RadenAbror/UserManagement/app/config"
	"github.com/RadenAbror/UserManagement/app/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// run database
	config.ConnectDB()

	// route link
	routes.Routes(router)

	// run server
	router.Run("localhost:1224")
}
