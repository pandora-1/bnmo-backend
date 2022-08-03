package database

import (
	"time"
)

type Transaction struct {
	Id          uint `json:"id" gorm:"primary_key, column:id"`
	Username       string             `json:"username,omitempty" gorm:"column:username"`
	Request float64             `json:"request,omitempty" validate:"required" gorm:"column:request"`
	IsSuccess       bool            `json:"isSuccess,omitempty" validate:"required" gorm:"column:isSuccess"`
	IsVerified       int            `json:"isVerified,omitempty" validate:"required" gorm:"column:isVerified"`
	CreatedAt   time.Time          `json:"created_at,omitempty" gorm:"column:createdAt"`
}