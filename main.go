package main

import(
	"net/http"
	"runtime"
	"api-gin-gen2/controllers"
	"github.com/gin-gonic/gin"
)

//Enable CORS (Cross-Origin Resource Sharing) - using middleware
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
	router.Use(CORS_Middleware())

	router.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", gin.H{
			"goVersion": runtime.Version(),
		})
	})
	auth := router.Group("/auth")
	{
		authController := new(controllers.AuthController)	
		auth.GET("/trial", func(c *gin.Context){
			c.HTML(http.StatusOK, "restricted.html", gin.H{"message": "welcome to trial page"})
		})
		auth.GET("/login", func(c *gin.Context){ c.HTML(http.StatusOK, "login.html", gin.H{}) })
		auth.GET("/register", func(c *gin.Context){ c.HTML(http.StatusOK, "register.html", gin.H{}) })
		auth.POST("/loginx", authController.Login);
		auth.POST("/registerx", authController.Register);
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
	- go get github.com/go-sql-driver/mysql

	reference ::
	https://github.com/evllnekrist/api-gin
	https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
	https://medium.com/@kiddy.xyz/tutorial-golang-rest-api-mysql-part-1-45cd9f4e75a6
	https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267

*/