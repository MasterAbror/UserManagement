package routes

import (
	"github.com/RadenAbror/UserManagement/app/controllers"
	"github.com/RadenAbror/UserManagement/app/helpers"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.POST("/user/create", controllers.CreateUser())
	router.POST("/user/auth", controllers.AuthUser())
	router.GET("/user/read/:userId", helpers.AuthorizationUser(), controllers.GetAUser())
	router.PUT("/user/update/:userId", helpers.AuthorizationUser(), controllers.EditAUser())
	router.DELETE("/user/delete/:userId", helpers.AuthorizationUser(), controllers.DeleteAUser())
	router.GET("/users", helpers.AuthorizationUser(), controllers.GetAllUsers())
	router.GET("/me", helpers.AuthorizationUser(), controllers.GetMe())
	router.GET("/user/logout", helpers.AuthorizationUser(), controllers.LogoutUser)
}
