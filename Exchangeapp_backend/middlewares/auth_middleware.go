package middlewares

import (
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户必须要登录或者注册, 有了 jwt 才可以创建对应的汇率内容。
//中间件是为了过滤路由而发明的一种机制, 也就是 http 请求来到时先经过中间件, 再到具体的处理函数。
//- 请求到到达我们定义的处理函数之前, 拦截请求并进行相应处理(比如: 权限验证, 数据过滤等), 这个可以类比为前置拦截器或前置过滤器;

// 验证jwt 存在和合规
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中返回值
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			// c.Abort 立即终止当前请求的处理流程
			// c.Abort 的作用是确保在验证失败时，请求不会继续执行后续的逻辑,注意是跳过后续的中间件和路由处理函数
			// 换句话说，c.Abort() 只会影响 Gin 的请求处理链，而不会影响当前函数的执行流程。
			// 在中间件中，如果验证失败，通常需要调用 c.Abort 终止请求
			c.Abort()
			return

		}

		username, err := utils.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// `c *gin.Context` -> `func1` -> `func2`(`c.Set(key,value)`) => `func3`(`c.Get(key)`)
		c.Set("username", username)

		// 调用下一个中间件或处理器
		c.Next()
	}
}
