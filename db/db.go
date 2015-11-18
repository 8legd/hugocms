package db

import (
	"github.com/jinzhu/gorm"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// QOR database
// TODO look at moving to some kind of file store / git repo api for better Hugo integration
var DB *gorm.DB

func init() {
	//if db, err := gorm.Open("mysql", "<username>:<password.@tcp(<host>:<port>)/qor_cms?charset=utf8&parseTime=True&loc=Local"); err != nil {
	if db, err := gorm.Open("sqlite3", "tmp/qor_cms.db"); err != nil {
		panic(err) //TODO more graceful!
	} else {
		DB = &db
	}
	DB.LogMode(true)
}
