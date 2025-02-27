package controllers

import (
	"exchangeapp/global"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func LikeArticle(c *gin.Context) {
	articleID := c.Param("id")

	// 设置Redis键， Redis中唯一标识某篇文章的点赞数
	// redis key 命名规范： key 单词与单词之间以 (:)分割, 如`user:userinfo`, `article:1:likes`
	likeKey := "article:" + articleID + ":likes"

	// Incr表示将 key 中储存的数字值增1
	if err := global.RedisDB.Incr(c, likeKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 一切没问题
	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

// 获取点赞数

func GetArticleLikes(c *gin.Context) {
	articleID := c.Param("id")

	likeKey := "article:" + articleID + ":likes"
	likes, err := global.RedisDB.Get(c, likeKey).Result()

	// 没有点赞数
	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}
