package controllers

import(
	"fmt"
	"net/http"
	"api-gin-gen2/models"
	"api-gin-gen2/helpers"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

var db_mysql = new(models.MySqlPool)//nama-package.nama-type-data-file-package
var helper_basic = new(helpers.EveBasicHelper)

func DisplayInfo(c *gin.Context, err error, info string) {
	if err != nil {
		if info != "" {
			info = "oops.. sorry, error occurring"
		}
		panic(err.Error())
	}
	c.HTML(http.StatusOK, "info.html", gin.H{"message": info})
}

// fmt.Println("")
// ----------------------------------------------------------------------------------- CRUD

type Users struct {
    id    int `json:"id"`
    name  string `json:"name"`
    username string `json:"username"`
    email string `json:"email"`
    password string `json:"password"`
    user_role int `json:"user_role"`
    user_status int `json:"user_status"`
    remember_token string `json:"remember_token"`
    signature_path string `json:"signature_path"`
    created_at string `json:"created_at"`
    updated_at string `json:"updated_at"`
    created_by string `json:"created_by"`
    updated_by string `json:"updated_by"`
}

func (ctrl AuthController) Register(c *gin.Context){
	db1 := db_mysql.MySQLConn() 

	r := c.Request
    if r.Method == "POST" {
		name 	:= r.FormValue("name")
		email 	:= r.FormValue("email")
		pwd 	:= r.FormValue("password")
		role 	:= r.FormValue("user_role")
		insertAct, err 	:= db1.Prepare("INSERT INTO gome_users (name,username,email,password,user_role,user_status,remember_token,signature_path,created_at,created_by,updated_at,updated_by) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
		resultAct, _	:= insertAct.Exec(name,email,email,pwd,role,1,"qwertyuiopasdfghjklzxcvbnm","signature.jpg","2019-08-28 00:00:00","default","2019-08-28 00:00:00","default")
		helper_basic.Panics(err)
		fmt.Println("REGISTER -- REGISTER -- REGISTER -- REGISTER -- REGISTER -- REGISTER -- REGISTER -- REGISTER")
		fmt.Println(resultAct)
	}
	defer db1.Close()
}

func (ctrl AuthController) Login(c *gin.Context){
	db1 := db_mysql.MySQLConn() 

	r := c.Request
    if r.Method == "POST" {
		email 	:= r.FormValue("email")
		// pwd 	:= r.FormValue("password")
		selectAct, err 	:= db1.Query("SELECT * FROM gome_users WHERE email = ?", email)

		user := Users{}
		res := []Users{}
		for selectAct.Next(){
		    err = selectAct.Scan(&user.id, &user.name, &user.username, &user.email, &user.password, &user.user_role, &user.user_status, 
		    	&user.remember_token, &user.signature_path, &user.created_at, &user.updated_at, &user.created_by, &user.updated_by)
			helper_basic.Panics(err)
        	res = append(res, user)
		}
		fmt.Println("LOGIN -- LOGIN -- LOGIN -- LOGIN -- LOGIN -- LOGIN -- LOGIN -- LOGIN")
		fmt.Println(res)
	}
	defer db1.Close()
}

