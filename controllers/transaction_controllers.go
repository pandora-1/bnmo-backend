package controllers

import (
	"github.com/gin-gonic/gin"
	"context"
	"time"
	"bnmo-backend/database"
	"bnmo-backend/responses"
	"fmt"
	"strconv"
    "net/http"
)

func CreateNewTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		username := c.PostForm("username")
		matauang := c.PostForm("matauang")
		hasil := c.PostForm("total")
		tambahuang := c.PostForm("tambahuang")
		total, _ := strconv.ParseFloat(hasil,64)
		hasilconvert, _ := ConvertUang(matauang,total)
		if tambahuang == "false" {
			hasilconvert *= -1
		}
		fmt.Println(matauang)
		newTransaction := database.Transaction{
			Username : username,
			Request : hasilconvert,
			IsSuccess : true,
			IsVerified : 1,
			CreatedAt : time.Now(),
		}
		database.DB.Create(&newTransaction)
		c.JSON(200, responses.UserResponse{Status: 200, Message: "Sukses menambah transaksi, silahkan menunggu verifikasi dari admin"})
		return
	}
}

func GetAllTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transaction []database.Transaction
		defer cancel()

		database.DB.Find(&transaction)
		c.JSON(http.StatusOK, gin.H{"data": transaction})
	}
}

func GetTransactionNotYetAccepted() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transaction []database.Transaction
		defer cancel()
		
		database.DB.Where("isVerified = ?", 1).Find(&transaction)

		c.JSON(http.StatusOK, gin.H{"data": transaction})
	}
}

func GetHistoryTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transaction []database.Transaction
		defer cancel()

		username := c.PostForm("username")

		database.DB.Where("username = ?", username).Find(&transaction)

		c.JSON(http.StatusOK, gin.H{"data": transaction})
	}
}

func ValidateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var transaction database.Transaction
		var user database.User
		defer cancel()
		id := c.PostForm("id")
		hasil := c.PostForm("accepted")
		accepted, _ := strconv.Atoi(hasil)
		if err := database.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
			c.JSON(400, gin.H{"message": "id tidak ditemukan"})
			return
		} else {
			database.DB.Model(&database.Transaction{}).Where("id = ?", id).Update("isVerified", accepted)
			if accepted == 2 {
				database.DB.Where("username = ?", transaction.Username).First(&user)
				saldo := user.Saldo + transaction.Request
				database.DB.Model(&database.User{}).Where("username = ?", user.Username).Update("saldo", saldo)
			}
			c.JSON(200, responses.UserResponse{Status: 200, Message: "Berhasil mengupdate data"})
			return
		}
	}
}