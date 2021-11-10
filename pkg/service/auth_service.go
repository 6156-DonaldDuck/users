package service

import (
	"encoding/json"
	"fmt"
	"github.com/6156-DonaldDuck/users/pkg/auth"
	"github.com/6156-DonaldDuck/users/pkg/config"
	"github.com/6156-DonaldDuck/users/pkg/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  auth.RedirectUrl,
		Scopes:       []string{auth.ScopeEmail, auth.ScopeUserProfile},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = ""
)

func init() {
	oauthConf.ClientID = config.Configuration.OAuth.ClientID
	oauthConf.ClientSecret = config.Configuration.OAuth.ClientSecret
	oauthStateString = config.Configuration.OAuth.OauthStateString
}

func BuildGoogleOAuthLoginURL() (string, error) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Errorf("[service.BuildGoogleOAuthLoginURL] error occurred while parsing auth url of endpoint: err=%v\n", err.Error())
		return "", err
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	log.Infof("[service.BuildGoogleOAuthLoginURL] successfully built google oauth login url: %s\n", url)
	return url, nil
}

func GoogleOAuthCallbackHandler(c *gin.Context, state, code string) (*oauth2.Token, error) {
	// check whether state matches
	if state != oauthStateString {
		err := fmt.Errorf("invalid oauth state, expected %s, got %s", oauthStateString, state)
		log.Errorf("[service.GoogleOAuthCallbackHandler] %v\n", err)
		return nil, err
	}

	// check whether code exists
	if code == "" {
		return nil, fmt.Errorf("code not found")
	} else {
		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			err = fmt.Errorf("exchange failed, err=%v\n", err)
			return nil, err
		}
		log.Infof("[service.GoogleOAuthCallbackHandler] access token=%v\n", token.AccessToken)
		log.Infof("[service.GoogleOAuthCallbackHandler] expiration time=%v\n", token.Expiry.String())
		log.Infof("[service.GoogleOAuthCallbackHandler] refresh token=%v\n" + token.RefreshToken)

		return token, nil
	}
}

func GetGoogleUserProfile(token *oauth2.Token) (*model.GoogleUserProfile, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Errorf("[service.GetGoogleUserProfile] err occurred while verifying token, err=%v\n", err)
		return nil, err
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[service.GetGoogleUserProfile] failed to read the response body, err=%v\n", err)
		return nil, err
	}

	userProfile := model.GoogleUserProfile{}
	err = json.Unmarshal(response, &userProfile)
	if err != nil {
		log.Errorf("[service.GetGoogleUserProfile] failed to unmarshal user profile, err=%v\n", err)
		return nil, err
	}
	return &userProfile, nil
}

func GetDBUserRelatedToGoogleUser(profile *model.GoogleUserProfile) (*model.User, error) {
	email := profile.Email
	return GetUserByEmail(email)
}

func CreateDBUserRelatedToGoogleUser(profile *model.GoogleUserProfile) (uint, error) {
	user := model.User{
		FirstName: profile.GivenName,
		LastName: profile.FamilyName,
		Email: profile.Email,
	}
	return CreateUser(user)
}