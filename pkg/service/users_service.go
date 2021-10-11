package service

import (
	"github.com/6156-DonaldDuck/users/pkg/db"
	"github.com/6156-DonaldDuck/users/pkg/model"
	log "github.com/sirupsen/logrus"
)

func ListUsers() ([]model.User, error) {
	var users []model.User
	result := db.DbConn.Find(&users)
	if result.Error != nil {
		log.Errorf("[service.ListUsers] error occurred while listing users, err=%v\n", result.Error)
	} else {
		log.Infof("[service.ListUsers] successfully listed users, rows affected = %v\n", result.RowsAffected)
	}
	return users, result.Error
}

func GetUserById(userId uint) (model.User, error) {
	user := model.User{}
	result := db.DbConn.First(&user, userId)
	if result.Error != nil {
		log.Errorf("[service.GetUserById] error occurred while getting user with id %v, err=%v\n", userId, result.Error)
	} else {
		log.Infof("[service.GetUserById] successfully got user with id %v, rows affected = %v\n", userId, result.RowsAffected)
	}
	return user, result.Error
}

func DeleteUserById(userId uint) (error) {
	user := model.User{}
	result := db.DbConn.Delete(&user, userId)
	if result.Error != nil {
		log.Errorf("[service.DeleteUserById] error occurred while deleting user with id %v, err=%v\n", userId, result.Error)
	} else {
		log.Infof("[service.DeleteUserById] successfully deleted user with id %v, rows affected = %v\n", userId, result.RowsAffected)
	}
	return result.Error
}

func CreateUser(user model.User) (uint, error) {
	result := db.DbConn.Create(&user)
	if result.Error != nil {
		log.Errorf("[service.CreateUser] error occurred while creating user, err=%v\n", result.Error)
	} else {
		log.Infof("[service.CreateUser] successfully created user with id %v, rows affected = %v\n", user.ID, result.RowsAffected)
	}
	return user.ID, result.Error
}

func UpdateUser(updateInfo model.User) (error){
	result := db.DbConn.Model(&updateInfo).Updates(updateInfo)
	if result.Error != nil {
		log.Errorf("[service.UpdateUser] error occurred while updating user, err=%v\n", result.Error)
	} else {
		log.Infof("[service.UpdateUser] successfully updated user with id %v, rows affected = %v\n", updateInfo.ID, result.RowsAffected)
	}
	return result.Error
}




