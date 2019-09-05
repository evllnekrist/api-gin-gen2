package main

import (
	"api-gin-gen2/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
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

	authController := new(controllers.AuthController) //taro diluar aja krn dipakai juga di group yang lain
	data := make(map[string]interface{})
	data["auth"] = authController.CheckAuth

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"message":               data,
			"ginBoilerplateVersion": "vxx",
			"goVersion":             runtime.Version(),
		})
	})

	auth := router.Group("/auth")
	{
		auth.GET("/", func(c *gin.Context) {
			data["info"] = "welcome to the page of Authentication"
			c.HTML(http.StatusOK, "restricted.html", gin.H{"message": data})
		})
		auth.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{"message": data}) })
		auth.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", gin.H{"message": data}) })
		auth.POST("/loginx", authController.Login)
		auth.POST("/registerx", authController.Register)
		auth.GET("/refresh", authController.NewToken)
		auth.GET("/logout", authController.Logout)
		auth.GET("/list", authController.List)
	}
	asm := router.Group("/asm")
	{
		asmController := new(controllers.AsmController)
		asm.GET("/", func(c *gin.Context) {
			data["info"] = "welcome to the page of Assembler"
			c.HTML(http.StatusOK, "restricted.html", gin.H{"message": data})
		})
		asm.GET("/compiler-phases/scanner", asmController.Scanner)
		asm.GET("/compiler-phases/parser/:act", asmController.Parser)
		// asm.GET("/compiler-phases/code-generator", asmController.CodeGenerator)
	}
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{"message": data})
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

_______________________________________________________________________________________


	reference ::

	- universal
	https://github.com/evllnekrist/api-gin
	https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
	https://medium.com/@kiddy.xyz/tutorial-golang-rest-api-mysql-part-1-45cd9f4e75a6
	https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
	https://kodingin.com/tutorial-gin-gonic-belajar-membuat-router-pada-gin-go/ <-- utk route khususnya yg ada parsing variable
	- auth
	https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
	- asm
	https://golang.org/doc/asm
	https://getstream.io/blog/how-a-go-program-compiles-down-to-machine-code/

*/
