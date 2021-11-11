package model

type GoogleLoginCallbackResponse struct {
	AccessToken string `json:"access_token"`
	UserId uint `json:"user_id"`
}