package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (s *server) handleCallback(w http.ResponseWriter, r *http.Request) {
	// validate request state
	if !s.authService.IsValidState(r.FormValue("state")) {
		errorJsonRespond(w, http.StatusBadRequest, errors.New("invalid state"))
		return
	}
	// exchange code for token
	var token , err = s.authService.ExchangeCode(r.FormValue("code"))
	if err != nil {
		errorJsonRespond(w, http.StatusInternalServerError, errors.New("failed to exchange code for token"))
		return
	}
	// get user info
	info, err := s.authService.GetInfoFromJWT(token.AccessToken)
	if err != nil {
		log.Printf("failed to get info from JWT: %s", err)
		errorJsonRespond(w, http.StatusInternalServerError, errors.New("failed to get info from JWT"))
		return
	}

	jsonInfo, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		log.Printf("failed to serialize claims to json: %s", err)
		errorJsonRespond(w, http.StatusInternalServerError, errors.New("failed to serialize claims to json"))
		return
	}

	httpInfo, err := s.authService.GetInfoFromHttp(info["sub"].(string), token.AccessToken)
	if err != nil {
		log.Printf("failed to get info from auth server: %s", err)
		errorJsonRespond(w, http.StatusInternalServerError, errors.New("failed to serialize claims to json"))
		return
	}

	// final response
	fmt.Fprintf(w, "User info from JWT: %s \n User public profile from http request: %s", jsonInfo, httpInfo)
}
