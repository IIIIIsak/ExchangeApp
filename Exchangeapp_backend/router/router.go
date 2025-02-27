package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	// 创建一个默认的gin路由引擎
	r := gin.Default()

	// 创建一个路由分组 "/api/auth/
	auth := r.Group("/api/auth")
	{
		// 当客户端发送POST请求到/api/auth/login 就会调用 controllers.login
		// 为了生效，需要再controllers 包中定义他们
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRate)
	// 调用中间件
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticlesById)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
	}
	return r
}
