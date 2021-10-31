package service

import (
	"errors"
	"github.com/6156-DonaldDuck/users/pkg/db"
	"github.com/6156-DonaldDuck/users/pkg/model"
	log "github.com/sirupsen/logrus"
	"context"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func ListUsers(offset int, limit int) ([]model.User, error) {
	var users []model.User
	result := db.DbConn.Limit(limit).Offset(offset).Find(&users)
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

func ListAddresses(offset int, limit int) ([]model.Address, error) {
	var addresses []model.Address
	result := db.DbConn.Limit(limit).Offset(offset).Find(&addresses)
	if result.Error != nil {
		log.Errorf("[service.ListAddresses] error occurred while listing address, err=%v\n", result.Error)
	} else {
		log.Infof("[service.ListAddresses] successfully listed address, rows affected = %v\n", result.RowsAffected)
	}
	return addresses, result.Error
}

func GetAddressById(addressId uint) (model.Address, error) {
	address := model.Address{}
	result := db.DbConn.First(&address, addressId)
	if result.Error != nil {
		log.Errorf("[service.GetAddressById] error occurred while getting address with id %v, err=%v\n", addressId, result.Error)
	} else {
		log.Infof("[service.GetAddressById] successfully got address with id %v, rows affected = %v\n", addressId, result.RowsAffected)
	}
	return address, result.Error
}

func CreateAddress(address model.Address) (uint, error) {
	err := VerifyUSStreetAddress(address)
	if err != nil {
		return 0, err
	} else {
		log.Infof("[service.CreateAddress] successfully verified address\n")
	}
	result := db.DbConn.Create(&address)
	if result.Error != nil {
		log.Errorf("[service.CreateAddress] error occurred while creating address, err=%v\n", result.Error)
	} else {
		log.Infof("[service.CreateAddress] successfully created address with id %v, rows affected = %v\n", address.ID, result.RowsAffected)
	}
	return address.ID, result.Error
}

func UpdateAddressById(updateInfo model.Address) (error){
	err := VerifyUSStreetAddress(updateInfo)
	if err != nil {
		return err
	} else {
		log.Infof("[service.CreateAddress] successfully verified address\n")
	}
	result := db.DbConn.Model(&updateInfo).Updates(updateInfo)
	if result.Error != nil {
		log.Errorf("[service.UpdateAddress] error occurred while updating address, err=%v\n", result.Error)
	} else {
		log.Infof("[service.UpdateAddress] successfully updated address with id %v, rows affected = %v\n", updateInfo.ID, result.RowsAffected)
	}
	return result.Error
}

func DeleteAddressById(addressId uint) (error) {
	address := model.Address{}
	result := db.DbConn.Delete(&address, addressId)
	if result.Error != nil {
		log.Errorf("[service.DeleteAddressById] error occurred while deleting address with id %v, err=%v\n", addressId, result.Error)
	} else {
		log.Infof("[service.DeleteAddressById] successfully deleted address with id %v, rows affected = %v\n", addressId, result.RowsAffected)
	}
	return result.Error
}

func GetAddressByUserId(userId uint) (model.Address, error) {
	address := model.Address{}
	user, err := GetUserById(userId)
	if err != nil {
		log.Errorf("[service.GetAddressByUserId] error occurred while getting user with id %v, err=%v\n", userId, err)
		return address, err
	}
	if user.AddressID == 0 {
		return address, errors.New("[service.GetAddressByUserId] user don't have address")
	}
	address, err = GetAddressById(user.AddressID)
	if err != nil {
		log.Errorf("[service.GetAddressByUserId] error occurred while getting address with id %v, err=%v\n", user.AddressID, err)
		return address, err
	}
	return address, err
}

func VerifyUSStreetAddress(address model.Address) (error) {
	client := wireup.BuildUSStreetAPIClient(
		wireup.SecretKeyCredential("2a2e8396-bfba-7363-7978-c189477b1291", "r3LK0jgf3N8HAsYU0ZqS"),
		// The appropriate license values to be used for your subscriptions
		// can be found on the Subscriptions page the account dashboard.
		// https://www.smartystreets.com/docs/cloud/licensing
		// wireup.WithLicenses("us-rooftop-geocoding-cloud"),
		// wireup.ViaProxy("https://my-proxy.my-company.com"), // uncomment this line to point to the specified proxy.
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
		// ...or maybe you want to supply your own http client:
		// wireup.WithHTTPClient(&http.Client{Timeout: time.Second * 30})
	)
	lookup1 := &street.Lookup{
		Street:        address.StreetName1,
		Street2:       address.StreetName2,
		Urbanization:  "", // Only applies to Puerto Rico addresses
		City:          address.City,
		State:         address.Region,
		ZIPCode:       address.PostalCode,
		MaxCandidates: 3,
		MatchStrategy: street.MatchStrict, 
	}

	batch := street.NewBatch()
	batch.Append(lookup1)
	if err := client.SendBatchWithContext(context.Background(), batch); err != nil {
		log.Errorf("[service.VerifyUSStreetAddress] error occurred while sending request to smartystreet\n")
		return err
	}

	for _, input := range batch.Records() {
		if len(input.Results) == 0 {
			log.Errorf("[service.VerifyUSStreetAddress] invalid input address\n")
			return errors.New("invalid input address")
		}
	}
	return nil
}