package service

import (
	"github.com/6156-DonaldDuck/users/pkg/db"
	"github.com/6156-DonaldDuck/users/pkg/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

func ListUsers(offset int, limit int) ([]model.User, int, error) {
	var users []model.User
	var totalCount int64
	result := db.DbConn.Limit(limit).Offset(offset).Find(&users)

	if result.Error != nil {
		log.Errorf("[service.ListUsers] error occurred while listing users, err=%v\n", result.Error)
	} else {
		log.Infof("[service.ListUsers] successfully listed users, rows affected = %v\n", result.RowsAffected)
	}
	db.DbConn.Model(model.User{}).Count(&totalCount)

	return users, int(totalCount), result.Error
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


func GetUserAddressById(userId uint) (model.UserAddress, error) {
	user := model.User{}
	address := model.Address{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(userId uint, user *model.User) {
		result := db.DbConn.First(&user, userId)
		if result.Error != nil {
			log.Errorf("[service.GetUserAddressById] error occurred while getting user with id %v, err=%v\n", userId, result.Error)
		} else {
			log.Infof("[service.GetUserAddressById] successfully got user with id %v, rows affected = %v\n", userId, result.RowsAffected)
		}
		wg.Done()
	}(userId, &user)

	go func(userId uint, address *model.Address) {
		result := db.DbConn.Where("user_id = ?", userId).First(&address)
		if result.Error != nil {
			log.Errorf("[service.GetAddressByUserId] error occurred while getting address by user id %v, err=%v\n", userId, result.Error)
		}
		wg.Done()
	}(userId, &address)
	
	wg.Wait()

	useraddress := model.UserAddress{
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email: user.Email,
		StreetName1: address.StreetName1,
		StreetName2: address.StreetName2,
		City: address.City,
		Region: address.Region,
		CountryCode: address.CountryCode,
		PostalCode: address.PostalCode,
	}

	return useraddress, nil
}


func GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}

	result := db.DbConn.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Errorf("[service.GetUserByEmail] error occurred while getting user with email %v, err=%v\n", email, result.Error)
		return nil, result.Error
	} else {
		log.Infof("[service.GetUserByEmail] successfully got user with email %v, rows affected = %v\n", email, result.RowsAffected)
		return &user, nil
	}
}

func DeleteUserById(userId uint) error {
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

func UpdateUser(updateInfo model.User) error {
	result := db.DbConn.Model(&updateInfo).Updates(updateInfo)
	if result.Error != nil {
		log.Errorf("[service.UpdateUser] error occurred while updating user, err=%v\n", result.Error)
	} else {
		log.Infof("[service.UpdateUser] successfully updated user with id %v, rows affected = %v\n", updateInfo.ID, result.RowsAffected)
	}
	return result.Error
}