package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

type User struct {
	gorm.Model
	FirstName  string
	LastName  string
	PhoneNumber  string
	Email  string	
}

func dbGetData() User {	
	dsn := "dbuser:dbuserdbuser@tcp(e6156-db.c890fe3d965k.us-east-2.rds.amazonaws.com:3306)/hw1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.Create(&User{FirstName: "Huaxuan", LastName: "Gao", PhoneNumber: "1111", Email: "a@b.com"})
	var user User
  	db.First(&user, "first_name = ?", "Huaxuan") // find product with code D42
	return user
}

func main() {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		var user = dbGetData()
		fmt.Println(user)
		c.JSONP(http.StatusOK, user)
	})
	r.Run(":8080")
}
	