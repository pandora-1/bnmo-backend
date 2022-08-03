package database

type User struct {
	Id          uint `json:"id" gorm:"primary_key, column:id"`
	Name        string             `json:"name,omitempty" validate:"required" gorm:"column:name"`
	Username       string             `json:"username,omitempty" gorm:"column:username"`
	Password string             `json:"password,omitempty" validate:"required" gorm:"column:password"`
	KTP       string            `json:"ktp,omitempty" validate:"required" gorm:"column:ktp"`
	Saldo float64 `json:"saldo" validate:"required" gorm:"column:saldo"`
	IsAccepted   int          `json:"isAccepted" gorm:"column:isAccepted"`
	IsAdmin int `json:"isAdmin" validate:"required" gorm:"column:isAdmin"`
}