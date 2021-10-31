package router

import (
	"errors"
	"fmt"
	docs "github.com/6156-DonaldDuck/users/docs"
	"github.com/6156-DonaldDuck/users/pkg/auth"
	"github.com/6156-DonaldDuck/users/pkg/config"
	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/6156-DonaldDuck/users/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)



func InitRouter() {
	r := gin.Default()
	r.Use(cors.Default()) // default allows all origin
	docs.SwaggerInfo.BasePath = config.Configuration.Mysql.Host

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

	// google oauth apis
	r.GET("/api/v1/login/google/url", GetGoogleLoginUrl)
	r.POST("/api/v1/login/google/callback", GoogleLoginCallback)
	r.GET("/api/v1/users/google/profile", GetGoogleUserProfile)

	r.Run(":" + config.Configuration.Port)
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
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	offset, errOffset := strconv.Atoi(offsetStr)
	limit, errLimit := strconv.Atoi(limitStr)
	if errOffset != nil {
		log.Errorf("[router.ListUsers] failed to parse offset %v, err=%v\n", offsetStr, errOffset)
		c.JSON(http.StatusBadRequest, "invalid offset")
		return
	}
	if errLimit != nil {
		log.Errorf("[router.ListUsers] failed to parse limit %v, err=%v\n", limitStr, errLimit)
		c.JSON(http.StatusBadRequest, "invalid limit")
		return
	}
	users, err := service.ListUsers(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, users)
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
	} else{
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
		c.JSON(http.StatusNoContent, "Successfully delete user with id "+ idStr)
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
func CreateUser(c *gin.Context){
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil{
		c.JSON(http.StatusBadRequest, err)
	}
	if user.ID != 0{
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

// @BasePath /api/v1

// @Summary List All Addresses
// @Schemes
// @Description List all addresses
// @Tags Addresses
// @Accept json
// @Produce json
// @Success 200 {json} addresses
// @Failure 500 internal server error
// @Router /addresses [get]
func ListAddresses(c *gin.Context) {
	offsetStr := c.Param("offset")
	limitStr := c.Param("limit")
	offset, errOffset := strconv.Atoi(offsetStr)
	limit, errLimit := strconv.Atoi(limitStr)
	if errOffset != nil {
		log.Errorf("[router.ListAddresses] failed to parse offset %v, err=%v\n", offsetStr, errOffset)
		c.JSON(http.StatusBadRequest, "invalid offset")
		return
	}
	if errLimit != nil {
		log.Errorf("[router.ListAddresses] failed to parse limit %v, err=%v\n", limitStr, errLimit)
		c.JSON(http.StatusBadRequest, "invalid limit")
		return
	}
	addresses, err := service.ListAddresses(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, addresses)
	}
}

// @Summary Get Address By Address Id
// @Schemes
// @Description Get addresses by addresses id
// @Tags Addresses
// @Accept json
// @Produce json
// @Param ID path int true "the id of a specfic addresses"
// @Success 200 {json} addresses
// @Failure 400 invalid addresses id
// @Router /addresses/{addresseId} [get]
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
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, err.Error())
		} else{
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusOK, address)
	}
}

// @Summary Create Address
// @Schemes
// @Description Create Address
// @Tags Addresses
// @Accept json
// @Produce json
// @Param street_number query string false "Street Number"
// @Param street_name_1 query string false "Street Name Line 1"
// @Param street_name_2 query string false "Street Name Line 2"
// @Param city query string false "City"
// @Param region query string false "Region"
// @Param country_code query string false "Country Code"
// @Param postal_code query string false "Postal Code"
// @Success 200 {json} address id
// @Failure 400 invalid address id
// @Router /addresses/ [post]
func CreateAddress(c *gin.Context){
	address := model.Address{}
	if err := c.ShouldBind(&address); err != nil{
		c.JSON(http.StatusBadRequest, err)
	}
	if address.ID != 0{
		_, err := service.GetUserById(address.ID)
		if err == nil {
			c.JSON(http.StatusUnprocessableEntity, "Duplicate key")
		}
	}
	addressId, err := service.CreateAddress(address)
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusCreated, addressId)
	}
}

// @Summary Update Address By Address Id
// @Schemes
// @Description Update address by address id
// @Tags Addresses
// @Accept json
// @Produce json
// @Param ID path int true "the id of a specfic address"
// @Success 200 {json} update successfully
// @Failure 400 invalid address id
// @Router /addresses/{articleId} [put]
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

// @Summary Delete Address By Address Id
// @Schemes
// @Description Delete address by address id
// @Tags Addresses
// @Accept json
// @Produce json
// @Param ID query int true "the id of a specfic address"
// @Success 200 {json} delete successfully
// @Failure 400 invalid address id
// @Router /addresses/ [delete]
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
		c.JSON(http.StatusNoContent, "Successfully delete address with id "+ idStr)
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
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, err.Error())
		} else{
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusOK, address)
	}
}

func GetGoogleLoginUrl(c *gin.Context) {
	loginUrl, err := service.BuildGoogleOAuthLoginURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, loginUrl)
	}
}

func GoogleLoginCallback(c *gin.Context) {
	params := make(map[string]string)
	if err := c.ShouldBind(&params); err != nil {
		log.Error(err)
		c.Error(err)
		return
	} else {
		log.Infof("params=%v\n", params)
	}
	state := params["state"]
	code := params["code"]

	log.Infof("state=%s, code=%s\n", state, code)

	token, err := service.GoogleOAuthCallbackHandler(c, state, code)
	if err != nil {
		err = fmt.Errorf("err while handling login callback, err=%v\n", err)
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	} else {
		// return the access token to the frontend
		c.String(http.StatusOK, token.AccessToken)
		// set the token to the local memory storage
		auth.TokenStoreInstance.SetToken(token.AccessToken, token)
	}
}

func GetGoogleUserProfile(c *gin.Context) {
	accessToken := c.Query(auth.AccessToken)
	if accessToken == "" {
		err := fmt.Errorf("access token should not be empty")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	token := auth.TokenStoreInstance.GetToken(accessToken)
	if token == nil {
		err := fmt.Errorf("token not found for access token=%s", accessToken)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userProfile, err := service.GetGoogleUserProfile(token)
	if err != nil {
		err = fmt.Errorf("failed to verify google oauth token, err=%v\n", err)
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	} else {
		c.JSON(http.StatusOK, userProfile)
	}
}