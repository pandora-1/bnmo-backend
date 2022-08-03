package controllers

import (
	"github.com/gin-gonic/gin"
	"context"
	"time"
	"bnmo-backend/database"
	"strconv"
	"bnmo-backend/responses"
	// "fmt"
    "net/http"
)

func CreateNewTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user1 database.User
		var user2 database.User
		username := c.PostForm("username")
		usernameTujuan := c.PostForm("usernameTujuan")
		hasil := c.PostForm("total")
		matauang := c.PostForm("matauang")
		total, _ := strconv.ParseFloat(hasil, 64)
		total, _ = ConvertUang(matauang,total)
		if err := database.DB.Where("username = ? AND isAccepted = ?", username, 2).First(&user1).Error; err != nil {
			c.JSON(400, responses.UserResponse{Status: 400, Message: "Username tidak ditemukan"})
			return
		}
		if err := database.DB.Where("username = ? AND isAccepted = ?", usernameTujuan, 2).First(&user2).Error; err != nil {
			c.JSON(400, responses.UserResponse{Status: 400, Message: "Username tujuan tidak ditemukan"})
			return
		}
		if user1.Saldo < total {
			c.JSON(400, responses.UserResponse{Status: 400, Message: "saldo tidak cukup"})
			return
		}

		newTransfer := database.Transfer{
			Username : username,
			UsernameTujuan : usernameTujuan,
			IsSuccess : true,
			Total : total * -1,
			CreatedAt: time.Now(),
		}

		newTransferTujuan := database.Transfer{
			Username : usernameTujuan,
			UsernameTujuan : username,
			IsSuccess : true,
			Total : total,
			CreatedAt: time.Now(),
		}
		database.DB.Create(&newTransfer)
		database.DB.Create(&newTransferTujuan)
		database.DB.Model(&database.User{}).Where("username = ?", username).Update("saldo", user1.Saldo - total)
		database.DB.Model(&database.User{}).Where("username = ?", usernameTujuan).Update("saldo", user2.Saldo + total)
		c.JSON(200, responses.UserResponse{Status: 200, Message: "Success post the data. Please close this tab"})
		return
	}
}

func GetAllTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transfer []database.Transfer
		defer cancel()

		database.DB.Find(&transfer)

		c.JSON(http.StatusOK, gin.H{"data": transfer})
	}
}

func GetHistoryTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transfer []database.Transfer
		defer cancel()

		username := c.PostForm("username")

		database.DB.Where("username = ?", username).Find(&transfer)

		c.JSON(http.StatusOK, gin.H{"data": transfer})
	}
}