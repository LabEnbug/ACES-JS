package mysql

import (
	"backend/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"time"
)

var DB *sql.DB

func InitDBPool() {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DefaultDatabase.User, config.DefaultDatabase.Pass, config.DefaultDatabase.Host, config.DefaultDatabase.Port,
		config.DefaultDatabase.Name))
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR:", err)
		}
		panic(err)
	}
	DB.SetMaxOpenConns(5)
	DB.SetMaxIdleConns(3)
	DB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err = DB.Ping()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR:", err)
		}
		panic(err)
	}
}

func CloseDBPool() {
	err := DB.Close()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return
	}
}
