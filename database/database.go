package database

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func NewDatabase() {
    USER := "" // Masukkan pengguna
    PASS := "" // Masukkan password
    HOST := "" // Masukkan host
    DBNAME := "" // Masukkan dbname

    URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, 
    HOST, DBNAME)
    db, err := gorm.Open(mysql.Open(URL))

    if err != nil {
        panic("Failed to connect to database!")

    }
    fmt.Println("Database connection established")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&Transaction{})
    db.AutoMigrate(&Transfer{})
    DB = db

}