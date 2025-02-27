package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `binding:"required"` // 必需项
	Content string `binding:"required"`
	Preview string `binding:"required"`
	Likes   int    `gorm:"default:0"` // 点赞数，默认为0
}
