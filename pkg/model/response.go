package model

type GoogleLoginCallbackResponse struct {
	AccessToken string `json:"access_token"`
	UserId uint `json:"user_id"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}