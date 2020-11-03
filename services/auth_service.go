package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	authentication "ccl/ccl-auth-api/core"
	"ccl/ccl-auth-api/models"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

//Login returns an user token if user login is successful
func Login(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	fmt.Println("authbackend", authBackend)

	success, dbUser := authBackend.Authenticate(requestUser)
	if success {
		token, err := authBackend.GenerateToken(dbUser.Username, dbUser.Role)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(models.TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

//RefreshToken generates a new JWT token
func RefreshToken(requestUser *models.User) []byte {
	authBackend := authentication.InitJWTAuthenticationBackend()

	token, err := authBackend.GenerateToken(requestUser.Username, requestUser.Role)
	if err != nil {
		panic(err)
	}
	response, err := json.Marshal(models.TokenAuthentication{token})
	if err != nil {
		panic(err)
	}
	return response
}

//Logout
func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}
