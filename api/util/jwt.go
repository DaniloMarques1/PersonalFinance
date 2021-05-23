package util

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/danilomarques1/personalfinance/api/dto"
	jwt "github.com/dgrijalva/jwt-go"
)

type Payload struct {
	UserId int64 `json:"userId"`
	jwt.StandardClaims
}

const EXPIRES_TIME = 86400

func NewToken(id int64) (string, error) {
	now := time.Now()
	claims := Payload{UserId: id, StandardClaims: jwt.StandardClaims{ExpiresAt: now.Unix() + EXPIRES_TIME}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	return tokenStr, err
}

func VerifyToken(tokenStr string) (int64, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		log.Printf("%v\n", err)
		return -1, false
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		log.Println("Not ok")
		return -1, false
	}

	if payload.ExpiresAt < time.Now().Unix() {
		log.Println("Already expired")
		return -1, false
	}

	return payload.UserId, true

}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authSlice := strings.Split(authHeader, " ")
		if len(authSlice) < 2 {
			RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Missing authorization token"})
			return
		} else {
			token := authSlice[1]
			log.Printf("token = %v\n", token)

			userId, valid := VerifyToken(token)
			if !valid {
				RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Invalid token"})
				return
			} else {
				log.Printf("User id = %v\n", userId)
				r.Header.Set("userId", strconv.Itoa(int(userId)))
				next.ServeHTTP(w, r)
			}
		}
	})
}
