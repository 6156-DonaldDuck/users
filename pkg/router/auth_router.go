package router

import (
	"fmt"
	"github.com/6156-DonaldDuck/users/pkg/auth"
	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/6156-DonaldDuck/users/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InitAuthRouters(r *gin.Engine) {
	// google oauth apis
	r.GET("/api/v1/login/google/url", GetGoogleLoginUrl)
	r.POST("/api/v1/login/google/callback", GoogleLoginCallback)
	r.GET("/api/v1/users/google/profile", GetGoogleUserProfile)
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
	}

	// create a user in the database for the google user if not exist
	googleProfile, err := service.GetGoogleUserProfile(token)
	if err != nil {
		err = fmt.Errorf("failed to get google profile, err=%v\n", err)
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var userId uint
	dbUser, _ := service.GetDBUserRelatedToGoogleUser(googleProfile)
	if dbUser == nil { // the Google user is not related to a db user yet
		//err = fmt.Errorf("Failed to login! User is not registered.")
		//log.Error(err)
		//c.JSON(http.StatusNotFound, err)
		//return

		userId, err = service.CreateDBUserRelatedToGoogleUser(googleProfile)
		if err != nil {
			err = fmt.Errorf("failed to create db user for google user, err=%v", err)
			log.Error(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		log.Infof("created db user %d for Google user with email %s\n", userId, googleProfile.Email)
	} else {
		userId = dbUser.ID
	}
	// return the access token and related user id to the frontend
	c.JSON(http.StatusOK, model.GoogleLoginCallbackResponse{
		AccessToken: token.AccessToken,
		UserId: userId,
	})
	// set the token to the local memory storage
	auth.TokenStoreInstance.SetToken(token.AccessToken, token)
}

func GetGoogleUserProfile(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	token := auth.TokenStoreInstance.GetToken(accessToken)
	userProfile, err := service.GetGoogleUserProfile(token)
	if err != nil {
		err = fmt.Errorf("failed to get google user profile, err=%v\n", err)
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	} else {
		c.JSON(http.StatusOK, userProfile)
	}
}
