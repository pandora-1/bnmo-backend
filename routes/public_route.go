package routes

import (
	"bnmo-backend/controllers"
	"github.com/gin-gonic/gin"
)

func PublicRoute(router *gin.Engine) {
	router.POST("/register", controllers.CreateUser())
	router.POST("/login", controllers.LoginUser())
	router.POST("/admin/validate/public", controllers.ValidateUser())
	router.POST("/admin/create", controllers.CreateAdmin())
}
