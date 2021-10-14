package router

import (
	"github.com/6156-DonaldDuck/users/pkg/config"
	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/6156-DonaldDuck/users/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)


func InitRouter() {
	r := gin.Default()
	// users
	r.GET( "/api/v1/users", ListUsers)
	r.GET( "/api/v1/users/:userId", GetUserById)
	r.POST( "/api/v1/users", CreateUser)
	r.PUT( "/api/v1/users/:userId", UpdateUserById)
	r.DELETE( "/api/v1/users/:userId", DeleteUserById)
	// addresses
	r.GET( "/api/v1/addresses", ListAddresses)
	r.GET( "/api/v1/addresses/:addressId", GetAddressById)
	r.POST( "/api/v1/addresses", CreateAddress)
	r.PUT( "/api/v1/addresses/:addressId", UpdateAddressById)
	r.DELETE( "/api/v1/addresses/:addressId", DeleteAddressById)
	// // user + address
	r.GET( "/api/v1/users/:userId/address", GetAddressByUserId)
	// r.POST( "/api/v1/users/:userId/address", SetAddressByUserId)
	// r.GET( "/api/v1/addresses/:addressId/users", ListUsersByAddressId)
	// not sure what this means
	// r.POST( "/api/v1/addresses/:id/users", ) 
	r.Run(":" + config.Configuration.Port)
}

func ListUsers(c *gin.Context) {
	users, err := service.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func GetUserById(c *gin.Context) {
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetUserById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	user, err := service.GetUserById(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUserById(c *gin.Context){
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.DeleteUserById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	err = service.DeleteUserById(uint(userId))
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "Successfully delete user with id "+ idStr)
	}
}

func CreateUser(c *gin.Context){
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil{
		c.JSON(http.StatusBadRequest, err)
	}
	userId, err := service.CreateUser(user)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, userId)
	}
}

func UpdateUserById(c *gin.Context){
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.UpdateUserById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	updateInfo := model.User{}
	if err := c.ShouldBind(&updateInfo); err != nil{
		c.JSON(http.StatusBadRequest, err)
	}
	updateInfo.ID = uint(userId)
	err = service.UpdateUser(updateInfo)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "update successfully")
	}
}

func ListAddresses(c *gin.Context) {
	users, err := service.ListAddresses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func GetAddressById(c *gin.Context) {
	idStr := c.Param("addressId")
	addressId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetAddressById] failed to parse address id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid address id")
		return
	}
	address, err := service.GetAddressById(uint(addressId))
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, address)
	}
}

func CreateAddress(c *gin.Context){
	address := model.Address{}
	if err := c.ShouldBind(&address); err != nil{
		c.JSON(http.StatusBadRequest, err)
	}
	addressId, err := service.CreateAddress(address)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, addressId)
	}
}

func UpdateAddressById(c *gin.Context){
	idStr := c.Param("addressId")
	addressId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.UpdateAddressById] failed to parse address id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid address id")
		return
	}
	updateInfo := model.Address{}
	if err := c.ShouldBind(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	updateInfo.ID = uint(addressId)
	err = service.UpdateAddressById(updateInfo)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, updateInfo)
	}
}

func DeleteAddressById(c *gin.Context){
	idStr := c.Param("addressId")
	addressId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.DeleteAddressById] failed to parse address id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid address id")
		return
	}
	err = service.DeleteAddressById(uint(addressId))
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "Successfully delete address with id "+ idStr)
	}
}

func GetAddressByUserId(c *gin.Context){
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetAddressByUserId] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	address, err := service.GetAddressByUserId(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, address)
	}
}