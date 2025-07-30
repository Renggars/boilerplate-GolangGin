package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v)/%v?%v", ENV.DB_USERNAME, ENV.DB_PASSWORD, ENV.DB_URL, ENV.DB_DATABASE, "charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta")
	fmt.Println(connectionString)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}
