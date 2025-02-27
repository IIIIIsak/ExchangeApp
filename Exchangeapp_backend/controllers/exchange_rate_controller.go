package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CreateExchangeRate(c *gin.Context) {
	var exchangeRate models.ExchangeRate
	if err := c.ShouldBindJSON(&exchangeRate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate.Date = time.Now()

	// AutoMigrate 自动迁移数据库表结构。
	//它的作用是根据传入的 Go 结构体（模型）自动创建或更新数据库表结构。
	//如果表不存在，GORM 会根据结构体定义创建表。
	//如果表已经存在，GORM 会根据结构体的变化（如新增字段、修改字段类型等）更新表结构。
	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AutoMigrate 后再 Create
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 完事了，返向客户端返回一个 HTTP 响应，表示数据创建成功
	c.JSON(http.StatusCreated, exchangeRate)
}

// 获得汇率信息
func GetExchangeRate(c *gin.Context) {
	// 切片存储
	var exchangeRate []models.ExchangeRate

	if err := global.Db.Find(&exchangeRate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 404
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 500
		}
		return
	}

	c.JSON(http.StatusOK, exchangeRate)

}
