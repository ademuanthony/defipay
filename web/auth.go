package web

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID     string `json:"userId"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

func CreateToken(userid string, authorized bool) (string, error) {
	expirationTime := time.Now().Add(360 * time.Minute)
	claims := &Claims{
		UserID:     userid,
		Authorized: authorized,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenSls(r events.APIGatewayProxyRequest) string {
	bearToken := r.Headers["Authorization"]
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
