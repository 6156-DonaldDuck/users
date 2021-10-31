package auth

import "golang.org/x/oauth2"

type TokenStore map[string]*oauth2.Token

var TokenStoreInstance TokenStore

func init() {
	TokenStoreInstance = make(map[string]*oauth2.Token)
}

func (t *TokenStore) GetToken(accessToken string) *oauth2.Token {
	return (*t)[accessToken]
}

func (t *TokenStore) SetToken(accessToken string, token *oauth2.Token) {
	(*t)[accessToken] = token
}