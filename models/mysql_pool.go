package models

import(
	"database/sql"
	"api-gin-gen2/helpers"

	_"github.com/go-sql-driver/mysql"
)

type MySqlPool struct{}

var helper_basic = new(helpers.EveBasicHelper)

func (db_mysql MySqlPool) MySQLConn() *sql.DB { //function itu harus capital
	dbDriver := "mysql"
	dbName := "gome"
	dbUser := "root"
	dbPass := ""

	db, err := sql.Open(dbDriver,dbUser+":"+dbPass+"@/"+dbName)
	helper_basic.Panics(err)
	return db
}