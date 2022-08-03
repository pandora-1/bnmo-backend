package controllers

import (
	"fmt"
	"os"
	"strconv"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	jwt "github.com/dgrijalva/jwt-go"
)

/* func GetAllCurrency() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
 		url := "https://api.apilayer.com/exchangerates_data/symbols"
		client := &http.Client {}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("apikey", "QAWsyB8QB8mg5EVd9XK1fkp5IBfMSnNk")

		if err != nil {
			fmt.Println(err)
		}
		res, err := client.Do(req)
			if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := ioutil.ReadAll(res.Body)
		var result map[string]interface{}
    	json.Unmarshal([]byte(body), &result)
		c.JSON(http.StatusOK, gin.H{"data": result})
		return
	}
} */

func ConvertUang(matauang string, total float64) (float64, error) {
	url := "https://api.apilayer.com/exchangerates_data/convert?to=IDR&from=" + matauang + "&amount=" + strconv.FormatFloat(total, 'f', 6, 64)

	client := &http.Client {}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "QAWsyB8QB8mg5EVd9XK1fkp5IBfMSnNk")

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
		if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	data := result["result"]
	fmt.Println(data)
	hasil, ok := data.(float64)
	if !ok {
		fmt.Println("ERROR")
	}
	return hasil, nil
}

func GenerateToken(user_id uint, is_admin int) (string, error) {
	token_lifespan,err := strconv.Atoi("1")
	if err != nil {
		return "",err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["is_admin"] = is_admin
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(c *gin.Context) (int, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 1, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		isAdmin, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["is_admin"]), 10, 32)
		if err != nil {
			return 1, err
		}
		return int(isAdmin), nil
	}
	return 1, err
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (uint, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

func JwtAuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, err := TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		if isAdmin != 0 {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

func JwtAuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, err := TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		if isAdmin != 1 {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
		return
	}
}