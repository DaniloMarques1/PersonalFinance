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

/*
	each time token expires i will use refresh token to request a new one
	and when returning a new token will also return a new refresh token.
	if refresh token is invalid, just return 401
*/

const EXPIRES_TIME_TOKEN = 86400
const EXPIRES_TIME_REFRESH_TOKEN = EXPIRES_TIME_TOKEN * 3

// returns token and refresh token
func NewToken(id int64) (string, string, error) {
	now := time.Now()
	tokenClaims := Payload{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Unix() + EXPIRES_TIME_TOKEN,
		},
	}
	token, err := generateToken(tokenClaims)
	if err != nil {
		return "", "", err
	}
	refreshTokenClaims := Payload{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Unix() + EXPIRES_TIME_REFRESH_TOKEN,
		},
	}

	refreshToken, err := generateToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func generateToken(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyToken(tokenStr string) (int64, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{},
		func(token *jwt.Token) (interface{}, error) {
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

			userId, valid := VerifyToken(token)
			if !valid {
				RespondJson(w, http.StatusUnauthorized, &dto.ErrorResponseDto{Message: "Invalid token"})
				return
			} else {
				r.Header.Set("userId", strconv.Itoa(int(userId)))
				next.ServeHTTP(w, r)
			}
		}
	})
}
