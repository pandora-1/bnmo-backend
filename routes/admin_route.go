package routes

import (
	"bnmo-backend/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
	protected := router.Group("/admin")
	protected.Use(controllers.JwtAuthAdmin())
	protected.GET("/users", controllers.GetAllUsers())
	protected.GET("/not-verified", controllers.GetUsersNotYetAccepted())
	protected.POST("/verify", controllers.ValidateUser())
	protected.POST("/transaction/verify", controllers.ValidateTransaction())
	protected.GET("/transaction/not-verified",controllers.GetTransactionNotYetAccepted())
}
