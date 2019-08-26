package main

import(
	"net/http"
	"runtime"
	"github.com/gin-gonic/gin"
)

func CORS_Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", gin.H{
			"goVersion": runtime.Version(),
		})
	})
	auth := router.Group("/auth")
	{
		auth.GET("/trial", func(c *gin.Context){
			c.HTML(http.StatusOK, "restricted.html", gin.H{"message": "welcome to trial page"})
		})
		auth.GET("/login", func(c *gin.Context){ c.HTML(http.StatusOK, "login.html", gin.H{}) })
		auth.GET("/register", func(c *gin.Context){ c.HTML(http.StatusOK, "register.html", gin.H{}) })
	}
	router.NoRoute( func(c * gin.Context){
		c.HTML(404, "404.html", gin.H{})
	})

	router.LoadHTMLGlob("./public/views/*")
	router.Static("/public", "./public") //serve static files
	router.Static("/assets", "./assets") //serve static files
	router.Run(":8080")
}

/* 
	prepare ::
	- go get -u github.com/gin-gonic/gin
*/