package router

import (
	"errors"
	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/6156-DonaldDuck/users/pkg/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	log "github.com/sirupsen/logrus"
)

func InitAddressRouters(r *gin.Engine) {
	// addresses
	r.GET("/api/v1/addresses", ListAddresses)
	r.GET("/api/v1/addresses/:addressId", GetAddressById)
	r.POST("/api/v1/addresses", CreateAddress)
	r.PUT("/api/v1/addresses/:addressId", UpdateAddressById)
	r.DELETE("/api/v1/addresses/:addressId", DeleteAddressById)
	// // user + address
	r.GET("/api/v1/users/:userId/address", GetAddressByUserId)
	// r.POST( "/api/v1/users/:userId/address", SetAddressByUserId)
	// r.GET( "/api/v1/addresses/:addressId/users", ListUsersByAddressId)
	// not sure what this means
	// r.POST( "/api/v1/addresses/:id/users", )
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
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, errPage := strconv.Atoi(pageStr)
	pageSize, errPageSize := strconv.Atoi(pageSizeStr)
	if errPage != nil {
		log.Errorf("[router.ListAddresses] failed to parse page %v, err=%v\n", pageStr, errPage)
		c.JSON(http.StatusBadRequest, "invalid offset")
		return
	}
	if errPageSize != nil {
		log.Errorf("[router.ListAddresses] failed to parse page size %v, err=%v\n", pageSize, errPageSize)
		c.JSON(http.StatusBadRequest, "invalid limit")
		return
	}
	addresses, total, err := service.ListAddresses((page - 1) * pageSize, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
	} else {
		c.JSON(http.StatusOK, model.ListAddressesResponse{
			Addresses: addresses,
			Total: total,
			Page: page,
			PageSize: pageSize,
		})
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
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
func CreateAddress(c *gin.Context) {
	address := model.Address{}
	if err := c.ShouldBind(&address); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if address.ID != 0 {
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
func UpdateAddressById(c *gin.Context) {
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
func DeleteAddressById(c *gin.Context) {
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
		c.JSON(http.StatusNoContent, "Successfully delete address with id "+idStr)
	}
}

func GetAddressByUserId(c *gin.Context) {
	idStr := c.Param("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetAddressByUserId] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	address, err := service.GetAddressByUserId(uint(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.JSON(http.StatusOK, address)
	}
}