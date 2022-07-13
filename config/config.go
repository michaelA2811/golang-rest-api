package config

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "admin12345"
	dbName := "testdb"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	return gormDB

}
