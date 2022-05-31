package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(username string, password string, dbname string, hostname string) {

	dbConnection, err := gorm.Open(mysql.Open(username+":"+password+"@tcp("+hostname+":3306)/"+dbname+"?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic("connection failed to the database ")
	}

	DB = dbConnection

	fmt.Println("Connect established Successfully...")

}
