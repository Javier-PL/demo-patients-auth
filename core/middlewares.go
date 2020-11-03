package authentication

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

//RequireTokenAuthentication proofs if the token from request is a valid token in order to grant access.
func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return authBackend.PublicKey, nil

	})

	if err == nil && token.Valid { //&& !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireAdminPermissions(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return authBackend.PublicKey, nil

	})

	//claims, done := extractClaims(token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if err == nil && role == "admin" {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}

}

func RequireEditorPermissions(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return authBackend.PublicKey, nil

	})

	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if err == nil && (role == "admin" || role == "editor") {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}

}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {

	authBackend := InitJWTAuthenticationBackend()

	token := jwt.New(jwt.SigningMethodRS512)
	hmacSecretString, _ := token.SignedString(authBackend.privateKey)

	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
