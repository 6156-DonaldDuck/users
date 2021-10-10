package db

import (
	"fmt"
	"github.com/6156-DonaldDuck/users/pkg/config"
	"github.com/6156-DonaldDuck/users/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DbConn *gorm.DB

func init() {
	host := config.Configuration.Mysql.Host
	port := config.Configuration.Mysql.Port
	username := config.Configuration.Mysql.Username
	password := config.Configuration.Mysql.Password
	databaseName := config.Configuration.Mysql.DatabaseName
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, databaseName)
	DbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("[db.init] error occurred while creating database connection, err=%v\n", err)
	}

	createTables()
}

func createTables() {
	if !DbConn.Migrator().HasTable("users") {
		log.Infof("[db.createTables] table users not found, creating new one\n")
		if err := DbConn.Migrator().CreateTable(&model.Comment{}); err != nil {
			log.Errorf("[db.createTables] error occurred while creating table users, err=%v\n", err)
		}

		// insert test data
		testUser := model.User{
			Model: gorm.Model{
				ID: 1,
				CreatedAt: time.Now(),
			},
			FirstName: "Tester",
			LastName: "A",
			PhoneNumber: "123-444-4321",
			Email: "a@b.com"
		}
		result := DbConn.Create(&testUser)
		if result.Error != nil {
			log.Errorf("[db.createTables] error occurred while inserting test data, err=%v\n", result.Error)
		} else {
			log.Infof("[db.createTables] successfully inserted test data, rows affected=%v\n", result.RowsAffected)
		}
	}
}