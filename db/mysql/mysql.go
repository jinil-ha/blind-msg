package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/golog"

	"github.com/jinil-ha/blind-msg/utils/config"
)

const (
	maxIdleConn        = 5
	maxOpenConn        = 32
	maxDeadlockRetry   = 5
	mysqlDeadlockError = 1213
)

var database *sql.DB

func init() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?timeout=5s&parseTime=true&loc=UTC&charset=utf8",
		config.GetString("mysql.username"),
		config.GetString("mysql.password"),
		config.GetString("mysql.host"),
		config.GetString("mysql.port"),
		config.GetString("mysql.dbname"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		golog.Fatalf("DB Connection failed : %s", err)
		panic("DB Connection failed!!")
	}

	database = db
}
