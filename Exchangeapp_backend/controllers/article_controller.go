package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// 缓存
var cacheKey = "articles"

func CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 迁移
	if err := global.Db.AutoMigrate(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 插入内容
	if err := global.Db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 当新增文章时，除了更新数据库，还要更新缓存中的对应数据，这样下次读取的时候，缓存中的数据就是最新的，不会导致新增文章看不到的情况
	// 可以在新增文章时，直接删除缓存中的旧数据;下次读取时, 由于缓存没有命中, 从数据库中读取最新数据并重新写入缓存
	if err := global.RedisDB.Del(c, cacheKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功创建了
	c.JSON(http.StatusCreated, article)

}

func GetArticles(c *gin.Context) {

	//获取缓存数据
	cachedData, err := global.RedisDB.Get(c, cacheKey).Result()

	// 如果Redis中没有找到对应的缓存即缓存未命中，那么先从数据库中获取文章数据
	if err == redis.Nil {
		// 切片 要拿到全部文章
		var articles []models.Article

		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// 收到数据后， 将数据序列化并存储在Redis缓存中，并设置过期时间; 同时还要返回给客户端
		articleJSON, err := json.Marshal(articles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := global.RedisDB.Set(c, cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, articles)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		// 缓存命中
		var articles []models.Article
		// 反序列化 为文章列表，返回给客户端
		if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, articles)
	}

}

func GetArticlesById(c *gin.Context) {
	id := c.Param("id")

	var article models.Article

	if err := global.Db.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, article)
}
