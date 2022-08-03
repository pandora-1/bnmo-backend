package database

import (
	"time"
)

type Transfer struct {
	Id          uint `json:"id" gorm:"primary_key, column:id"`
	Username       string             `json:"username,omitempty" gorm:"column:username"`
	UsernameTujuan string             `json:"usernameTujuan,omitempty" validate:"required" gorm:"column:usernameTujuan"`
	IsSuccess       bool            `json:"isSuccess,omitempty" validate:"required" gorm:"column:isSuccess"`
	Total float64 `json:"total,omitempty" validate:"required" gorm:"column:total"`
	CreatedAt   time.Time          `json:"created_at,omitempty" gorm:"column:createdAt"`
}