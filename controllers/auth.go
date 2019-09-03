package controllers

import( 
	// "fmt"
	"net/http"
	"api-gin-gen2/models"
	"api-gin-gen2/helpers"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

var db_mysql = new(models.MySqlPool)//nama-package.nama-type-data-file-package
var helper_basic = new(helpers.EveBasicHelper)

// interface{} <---------- interface bisa untuk tampung semua tipe data
// fmt.Println("") || "fmt"
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
    	data := make(map[string]interface{})
		name 	:= r.FormValue("name")
		email 	:= r.FormValue("email")
		pwd 	:= r.FormValue("password")
		role 	:= r.FormValue("user_role")
		err_str := make(map[int]interface{})

		insertAct, err 	:= db1.Prepare("INSERT INTO gome_users (name,username,email,password,user_role,user_status,remember_token,signature_path,created_at,created_by,updated_at,updated_by) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
		resultAct, err2	:= insertAct.Exec(name,email,email,pwd,role,1,"qwertyuiopasdfghjklzxcvbnm","signature.jpg","2019-08-28 00:00:00","default","2019-08-28 00:00:00","default")
		err_str[0] = helper_basic.Panics(err)
		err_str[1] = helper_basic.Panics(err2)

    	data["error"] = err_str
    	data["result"] = resultAct
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
		
	}
	defer db1.Close()
}

func (ctrl AuthController) Login(c *gin.Context){
	db1 := db_mysql.MySQLConn() 

	r := c.Request
    if r.Method == "POST" {
    	data := make(map[string]interface{})
		user := Users{}

		email 	:= r.FormValue("email")
		pwd 	:= r.FormValue("password")

		err := db1.QueryRow("SELECT * FROM gome_users WHERE email = ? LIMIT 1", email).Scan(&user.id, &user.name, 
			&user.username, &user.email, &user.password, &user.user_role, &user.user_status, &user.remember_token, 
			&user.signature_path, &user.created_at, &user.updated_at, &user.created_by, &user.updated_by)
		err_str := helper_basic.Panics(err)
		if(err_str == nil && user.password != pwd){
			data["error"] = "password not match";
			user = Users{}
		}else{
    		data["error"] = err_str
		}

    	data["result"] = user
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
	}
	defer db1.Close()
}


func (ctrl AuthController) List(c *gin.Context){
	db1 := db_mysql.MySQLConn() 

	r := c.Request
    if r.Method == "POST" {
    	data := make(map[string]interface{})
		user := Users{}
		resultAct := []Users{}

		selectAct, err 	:= db1.Query("SELECT * FROM gome_users")

		for selectAct.Next(){
		    err = selectAct.Scan(&user.id, &user.name, &user.username, &user.email, &user.password, &user.user_role, &user.user_status, 
		    	&user.remember_token, &user.signature_path, &user.created_at, &user.updated_at, &user.created_by, &user.updated_by)
			helper_basic.Panics(err)
        	resultAct = append(resultAct, user)
		}

    	data["result"] = resultAct
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
	}
	defer db1.Close()
}

