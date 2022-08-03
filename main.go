package main

import (
  "bnmo-backend/database"
  "bnmo-backend/routes"
  "github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())
	database.NewDatabase() //new database connection
  var user database.User
  if err := database.DB.Where("username = ?", "admin").First(&user).Error; err == nil {

  } else {
    newUser := database.User{
      Name:        "admin",
      KTP:       "nothing",
      Username: "admin",
      Saldo: 0,
      Password: "admin",
      IsAccepted: 2,
      IsAdmin: 1,
    }
    database.DB.Create(&newUser)
  }

	routes.UserRoute(router)
  routes.AdminRoute(router)
  routes.PublicRoute(router)

	router.Run(":8000")
}