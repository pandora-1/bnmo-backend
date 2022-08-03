package controllers

import (
	"github.com/gin-gonic/gin"
	"context"
	"time"
	"bnmo-backend/database"
	"encoding/base64"
	"bnmo-backend/responses"
	"fmt"
    "net/http"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user database.User
		nama := c.PostForm("name")
		username := c.PostForm("username")
		password := c.PostForm("password")
		file, err := c.FormFile("ktp")

		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}
		if err := database.DB.Where("username = ? AND isAccepted = ?", username,2).First(&user).Error; err == nil {
			c.JSON(400, responses.UserResponse{Status: 400, Message: "Username sudah terdaftar"})
			return
		}
		fmt.Println(err)
		var base64Encoding string
		content, _ := file.Open()
		byteContainer, err := ioutil.ReadAll(content)
		mimeType := http.DetectContentType(byteContainer)
		switch mimeType {
		case "image/jpeg":
			base64Encoding += "data:image/jpeg;base64,"
		case "image/png":
			base64Encoding += "data:image/png;base64,"
		}
		base64Encoding += toBase64(byteContainer)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
		if err != nil {
			return
		}
		result := string(hashedPassword)

		newUser := database.User{
			Name:        nama,
			KTP:       base64Encoding,
			Username: username,
			Saldo: 0,
			Password: result,
			IsAccepted: 1,
			IsAdmin: 1,
		}
		database.DB.Create(&newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Terdapat kesalahan sistem"})
			return
		}
		c.JSON(200, responses.UserResponse{Status: 200, Message: "Sukses menambah pengguna, silahkan menunggu verifikasi dari admin"})
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user database.User
		nama := c.PostForm("name")
		username := c.PostForm("username")
		password := c.PostForm("password")
		file, err := c.FormFile("ktp")
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}
		err = database.DB.Where("username = ?", username).First(&user).Error
		if err == nil {
			c.JSON(400, responses.UserResponse{Status: 400, Message: "Username sudah terdaftar"})
			return
		}
		var base64Encoding string
		content, _ := file.Open()
		byteContainer, err := ioutil.ReadAll(content)
		mimeType := http.DetectContentType(byteContainer)
		switch mimeType {
		case "image/jpeg":
			base64Encoding += "data:image/jpeg;base64,"
		case "image/png":
			base64Encoding += "data:image/png;base64,"
		}
		base64Encoding += toBase64(byteContainer)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
		if err != nil {
			return
		}
		result := string(hashedPassword)

		newUser := database.User{
			Name:        nama,
			KTP:       base64Encoding,
			Username: username,
			Saldo: 0,
			Password: result,
			IsAccepted: 1,
			IsAdmin: 0,
		}
		database.DB.Create(&newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Terdapat kesalahan sistem"})
			return
		}
		c.JSON(200, responses.UserResponse{Status: 200, Message: "Sukses menambah pengguna, silahkan menunggu verifikasi dari admin"})
		return 
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []database.User
		defer cancel()

		database.DB.Where("isAdmin = ?", 0).Find(&users)

		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func GetUsersNotYetAccepted() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []database.User
		defer cancel()
		
		database.DB.Where("isAccepted = ?", 1).Find(&users)

		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func GetUsersById() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users database.User
		defer cancel()
		id := c.PostForm("id")
		fmt.Println(id)
		database.DB.Where("id = ?", id).Find(&users)

		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users database.User
		defer cancel()
		username := c.PostForm("username")
		password := c.PostForm("password")
		if err := database.DB.Where("username = ? AND isAccepted = ?", username, 2).Find(&users).Error; err != nil {
			c.JSON(400, gin.H{"message": "username tidak ditemukan"})
		} else {
			fmt.Println(users.Password)
			err = VerifyPassword(password, users.Password)
			if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				c.JSON(400, responses.UserResponse{Status: 400, Message: "Password salah"})
				return
			}
			if users.IsAccepted == 2 {		
				token,err := GenerateToken(users.Id, users.IsAdmin)
		
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "berhasil login", "data": token, "isAdmin": users.IsAdmin, "id": users.Id, "username":users.Username})
				return
			} else if users.IsAccepted == 1{
				c.JSON(400, responses.UserResponse{Status: 400, Message: "Belum disetujui"})
				return
			} else if users.IsAccepted == 3 {
				c.JSON(400, responses.UserResponse{Status: 400, Message: "Permohonan ditolak"})
				return
			}
		}
	}
}

func ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users database.User
		// var parameter MasukanValidate
		defer cancel()
		username := c.PostForm("username")
		accepted := c.PostForm("accepted")
		if err := database.DB.Where("username = ?", username).First(&users).Error; err != nil {
			c.JSON(400, gin.H{"message": "username tidak ditemukan"})
		} else {
			database.DB.Model(&database.User{}).Where("username = ?", username).Update("isAccepted", accepted)
			c.JSON(200, gin.H{"message": "behasil diupdate"})
		}
	}
}

func VerifyPassword(password,hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}