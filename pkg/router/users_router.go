package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/6156-DonaldDuck/users/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitUserRouters(r *gin.Engine) {
	r.GET("/api/v1/users", ListUsers)
	r.GET("/api/v1/users/:userId", GetUserById)
	r.POST("/api/v1/users", CreateUser)
	r.PUT("/api/v1/users/:userId", UpdateUserById)
	r.DELETE("/api/v1/users/:userId", DeleteUserById)
}

// @BasePath /api/v1

// @Summary List All Users
// @Schemes
// @Description List all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {json} users
// @Failure 500 internal server error
// @Router /users [get]
func ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, errPage := strconv.Atoi(pageStr)
	pageSize, errPageSize := strconv.Atoi(pageSizeStr)
	if errPage != nil {
		log.Errorf("[router.ListUsers] failed to parse page %v, err=%v\n", pageStr, errPage)
		c.JSON(http.StatusBadRequest, "invalid offset")
		return
	}
	if errPageSize != nil {
		log.Errorf("[router.ListUsers] failed to parse page size %v, err=%v\n", pageSize, errPageSize)
		c.JSON(http.StatusBadRequest, "invalid limit")
		return
	}
	users, total, err := service.ListUsers((page - 1) * pageSize, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, model.ListUsersResponse{
			Users: users,
			Total: total,
			Page: page,
			PageSize: pageSize,
		})
	}
}

// @Summary Get User By User Id
// @Schemes
// @Description Get user by user id
// @Tags Users
// @Accept json
// @Produce json
// @Param ID path int true "the id of a specfic user"
// @Success 200 {json} user
// @Failure 400 invalid user id
// @Router /users/{userId} [get]
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Error(err)
		}
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// @Summary Delete User By User Id
// @Schemes
// @Description Delete user by user id
// @Tags Users
// @Accept json
// @Produce json
// @Param ID query int true "the id of a specfic user"
// @Success 200 {json} delete successfully
// @Failure 400 invalid user id
// @Router /users/ [delete]
func DeleteUserById(c *gin.Context) {
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
		c.JSON(http.StatusNoContent, "Successfully delete user with id "+idStr)
	}
}

// @Summary Create User
// @Schemes
// @Description Create User
// @Tags Users
// @Accept json
// @Produce json
// @Param first_name query string false "First Name"
// @Param last_name query string false "Last Name"
// @Param phone_number query string false "Phone Number"
// @Param email query string false "Email"
// @Param address_id query int false "Address ID"
// @Success 200 {json} user id
// @Failure 400 invalid user id
// @Router /users/ [post]
func CreateUser(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if user.ID != 0 {
		_, err := service.GetUserById(user.ID)
		if err == nil {
			c.JSON(http.StatusUnprocessableEntity, "Duplicate key")
		}
	}
	userId, err := service.CreateUser(user)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusCreated, userId)
	}
}

// @Summary Update User By User Id
// @Schemes
// @Description Update user by user id
// @Tags Users
// @Accept json
// @Produce json
// @Param ID path int true "the id of a specfic user"
// @Success 200 {json} update successfully
// @Failure 400 invalid user id
// @Router /users/{userId} [put]
func UpdateUserById(c *gin.Context) {
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.UpdateUserById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	updateInfo := model.User{}
	if err := c.ShouldBind(&updateInfo); err != nil {
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