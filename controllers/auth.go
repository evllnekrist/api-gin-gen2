package controllers

import (
	"api-gin-gen2/helpers"
	"api-gin-gen2/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AuthController struct{}

// interface{} <---------- interface bisa untuk tampung semua tipe data
// fmt.Println("") || "fmt"
var db_mysql = new(models.MySqlPool) //nama-package.nama-type-data-file-package
var helper_basic = new(helpers.EveBasicHelper)
var jwtKey = []byte("my_secret_key")

type Users struct {
	id             int    `json:"id"`
	name           string `json:"name"`
	username       string `json:"username"`
	email          string `json:"email"`
	password       string `json:"password"`
	user_role      int    `json:"user_role"`
	user_status    int    `json:"user_status"`
	remember_token string `json:"remember_token"`
	signature_path string `json:"signature_path"`
	created_at     string `json:"created_at"`
	updated_at     string `json:"updated_at"`
	created_by     string `json:"created_by"`
	updated_by     string `json:"updated_by"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claim struct {
	username string `json:"username"`
	jwt.StandardClaims
}

// ----------------------------------------------------------------------------------- CRUD

func CheckAuth(c *gin.Context) int {
	// & returns the memory address of the following variable.
	// * (asterisk) returns the value of the following variable (which should hold the memory address of a variable,
	// unless you want to get weird output and possibly problems because you're accessing your computer's RAM)
	w := c.Writer
	r := c.Request
	data := make(map[string]interface{})

	//----------------------------checking authentication START------------------------untuk dipindah ke fungsi lain
	cookie, err_token := r.Cookie("gome_auth")

	if err_token != nil {
		if err_token == http.ErrNoCookie { //if the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			data["error"] = "TOKEN :: status unauthorized :: you not login yet or your session just expired"
			c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
			return 0
		} //for any other type of error, return bad request status
		w.WriteHeader(http.StatusBadRequest)
		data["error"] = "TOKEN :: status bad request :: " + err_token.Error() // kalau err_token, data typenya error kalau yg ada fx Error nya itu sdh jd string
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
		return 0
	}

	tokenString := cookie.Value
	claim := &Claim{}

	//parse the JWT string and store the result in `claims`.
	//note that we are passing the key in this method as well. This method will return an error
	//if the token is invalid (if it has expired according to the expiry time we set on sign in),
	//or if the signature does not match
	token, err_token := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err_token != nil {
		if err_token == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			data["error"] = "TOKEN :: status unauthorized :: signature invalid"
			c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
			return 0
		}
		w.WriteHeader(http.StatusBadRequest)
		data["error"] = "TOKEN :: status bad request :: " + err_token.Error()
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
		return 0
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		data["error"] = "TOKEN :: status unauthorized :: token invalid"
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
		return 0
	}

	return 1
}

func (ctrl AuthController) Register(c *gin.Context) {
	data := make(map[string]interface{})
	data["auth"] = CheckAuth(c)
	if data["auth"] == 0 { //kalau tdk terotentifikasi, tdk lanjut
		return
	}

	db1 := db_mysql.MySQLConn()
	r := c.Request
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		pwd := r.FormValue("password")
		role := r.FormValue("user_role")
		err_str := make(map[int]interface{})

		insertAct, err := db1.Prepare("INSERT INTO gome_users (name,username,email,password,user_role,user_status,remember_token,signature_path,created_at,created_by,updated_at,updated_by) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
		resultAct, err2 := insertAct.Exec(name, email, email, pwd, role, 1, "qwertyuiopasdfghjklzxcvbnm", "signature.jpg", "2019-08-28 00:00:00", "default", "2019-08-28 00:00:00", "default")
		err_str[0] = helper_basic.Panics(err)
		err_str[1] = helper_basic.Panics(err2)

		data["error"] = err_str
		data["user"] = resultAct
		data["result"] = "New user registered"
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})

	}
	defer db1.Close()
}

func (ctrl AuthController) Login(c *gin.Context) {

	db1 := db_mysql.MySQLConn()
	w := c.Writer
	r := c.Request
	if r.Method == "POST" {
		data := make(map[string]interface{})
		user := Users{}

		email := r.FormValue("email")
		pwd := r.FormValue("password")

		err := db1.QueryRow("SELECT * FROM gome_users WHERE email = ? LIMIT 1", email).Scan(&user.id, &user.name,
			&user.username, &user.email, &user.password, &user.user_role, &user.user_status, &user.remember_token,
			&user.signature_path, &user.created_at, &user.updated_at, &user.created_by, &user.updated_by)
		err_str := helper_basic.Panics(err)

		if err_str == nil && user.password != pwd {
			data["error"] = "password not match"
			user = Users{}
		} else if err_str == nil && user.password == pwd { //if no problem & password is match
			loc, _ := time.LoadLocation("Asia/Jakarta")
			expirationTime := time.Now().In(loc).Add(5 * time.Minute) //set expirasi to 5 menit
			claim := &Claim{
				username: user.username,
				StandardClaims: jwt.StandardClaims{ //in JWT, the expired time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenString, err2 := token.SignedString(jwtKey)
			if err2 != nil {
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				data["error"] = "TOKEN :: status internal server error :: cant make your session token"
				c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
				return
			}
			//finally, set client cookie for 'token' as the JWT just generated
			http.SetCookie(w, &http.Cookie{
				Name:    "gome_auth",
				Value:   tokenString,
				Expires: expirationTime,
			})

		} else {
			data["error"] = err_str
		}

		data["auth"] = 1
		data["user"] = user
		data["result"] = "You now logged"
		c.HTML(http.StatusOK, "info.html", gin.H{"message": data})
	}
	defer db1.Close()
}

func (ctrl AuthController) List(c *gin.Context) { //untested yet

	data := make(map[string]interface{})
	data["auth"] = CheckAuth(c)
	if data["auth"] == 0 { //kalau tdk terotentifikasi, tdk lanjut
		return
	}

	db1 := db_mysql.MySQLConn()
	r := c.Request
	if r.Method == "GET" {
		user := Users{}
		resultAct := []Users{}

		selectAct, err := db1.Query("SELECT * FROM gome_users")

		for selectAct.Next() {
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

func (ctrl AuthController) NewToken(c *gin.Context) {
	fmt.Println("FUNCTION NOT READY YET")
}

func (ctrl AuthController) Logout(c *gin.Context) {

	w := c.Writer
	// r := c.Request
	loc, _ := time.LoadLocation("Asia/Jakarta")
	expirationTime := time.Now().In(loc)

	http.SetCookie(w, &http.Cookie{
		Name:    "gome_auth",
		Value:   "",
		Expires: expirationTime,
	})
}
