package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type AuthService struct {
	oAuthConfig *oauth2.Config
	randomState string
	jwtAccess   string
	userInfoURL func(userId string) string
}

func NewAuthService() *AuthService {
	return &AuthService{
		oAuthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:5000/callback",
			ClientID:     "310",
			ClientSecret: "qOorjUDbg5OLNi7nNCOQrrrbOHA6KgLV",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://login.xsolla.com/api/oauth2/login",
				TokenURL:  "https://login.xsolla.com/api/oauth2/token",
				AuthStyle: 0,
			},
		},
		randomState: "myrandom",
		jwtAccess:   "QtGwsOOJHxepeqrJ7IJaMSZZFrwIyYp4Xz1JK2koDgDJIqxDJ4bm94iHNkCcyUQp",
		userInfoURL: func(userId string) string {
			return fmt.Sprintf("https://login.xsolla.com/api/users/%s/public", userId)
		},
	}
}

func (s *AuthService) IsValidState(state string) bool {
	return state == s.randomState
}

func (s *AuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	var token, err = s.oAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		err = fmt.Errorf("failed to exchange code for token: %s", err)
		return nil, err
	}
	return token, nil
}

func (s *AuthService) GetInfoFromJWT(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtAccess), nil
	}); err != nil {
		// sometimes tokens come signed with sign date parameter greater than current system time
		// maybe because system time and auth server time are not synced
		if err.Error() == "Token used before issued" {
			return claims, nil
		}
		err = fmt.Errorf("failed to parse jwt: %s", err)
		return nil, err
	}
	return claims, nil
}

func (s *AuthService) GetInfoFromHttp(userId string, token string) (string, error) {
	url := s.userInfoURL(userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("could not create get request : %s", err.Error())
		return "", err
	}
	req.Header.Add("authorization", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("could not get response : %s", err.Error())
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("could not parse response : %s", err.Error())
		return "", err
	}

	// prettify response body string
	var jsonStruct interface{}
	json.Unmarshal(body, &jsonStruct)
	jsonBody, _ := json.MarshalIndent(jsonStruct, "", "	")
	return string(jsonBody), nil
}
