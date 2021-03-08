package services

import "golang.org/x/oauth2"

type AuthServiceProvider interface {
	IsValidState(state string) bool
	ExchangeCode(code string) (*oauth2.Token, error)
	GetInfoFromJWT(token string) (map[string]interface{}, error)
	GetInfoFromHttp(userId string, token string) (string, error)
}
