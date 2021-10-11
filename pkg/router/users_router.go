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
	r.GET("/users", ListUsers)
	r.GET("/users/:userId", GetUserById)
	r.DELETE("/users/:userId", DeleteUserById)
	r.POST("/users", CreateUser)
	r.POST("/users/:userId", UpdateUserById)
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
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUserById(c *gin.Context){
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetUserById] failed to parse user id %v, err=%v\n", idStr, err)
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
		log.Errorf("[router.GetUserById] failed to parse user id %v, err=%v\n", idStr, err)
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



