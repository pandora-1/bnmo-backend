package routes

import (
	"bnmo-backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	protected := router.Group("/user")
	protected.Use(controllers.JwtAuthUser())
	protected.POST("/id", controllers.GetUsersById())
	protected.POST("/transfer", controllers.CreateNewTransfer())
	protected.GET("/transfer", controllers.GetAllTransfer())
	protected.POST("/transaction", controllers.CreateNewTransaction())
	protected.GET("/transaction", controllers.GetAllTransaction())
	protected.POST("/history/transaction", controllers.GetHistoryTransaction())
	protected.POST("/history/transfer", controllers.GetHistoryTransfer())
}
