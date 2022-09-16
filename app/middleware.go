package app

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
)

// CacheControl creates a new middleware to set the HTTP response header with
// "Cache-Control: max-age=maxAge" where maxAge is in seconds.
func CacheControl(maxAge int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(maxAge, 10))
			next.ServeHTTP(w, r)
		})
	}
}

func getClaims(r *http.Request) (Claims, error) {
	claims := Claims{}

	tknStr := ExtractToken(r)
	if tknStr == "" {
		log.Error("no token")
		return claims, errors.New("missing auth token")
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, errors.New("invalid auth token")
		}
		return claims, errors.New("error in processing")
	}
	if !tkn.Valid {
		return claims, errors.New("invalid auth token")
	}

	return claims, nil
}

func getClaimsSls(r events.APIGatewayProxyRequest) (Claims, error) {
	claims := Claims{}

	tknStr := ExtractTokenSls(r)
	if tknStr == "" {
		log.Error("no token")
		return claims, errors.New("missing auth token")
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, errors.New("invalid auth token")
		}
		return claims, errors.New("error in processing")
	}
	if !tkn.Valid {
		return claims, errors.New("invalid auth token")
	}

	return claims, nil
}

func (s Module) GetUserIDTokenCtx(r *http.Request) string {
	// Initialize a new instance of `Claims`
	claims, err := getClaims(r)
	if err != nil {
		return ""
	}

	if !claims.Authorized {
		return ""
	}

	return claims.UserID
}

func (s Module) GetUserIDTokenCtxSls(r events.APIGatewayProxyRequest) string {
	claims, err := getClaimsSls(r)
	if err != nil {
		return ""
	}

	if !claims.Authorized {
		return ""
	}

	return claims.UserID
}

func (s Module) GetUserIDTokenUnAuthCtxSls(r events.APIGatewayProxyRequest) string {
	// Initialize a new instance of `Claims`
	claims, err := getClaimsSls(r)
	if err != nil {
		return ""
	}

	return claims.UserID
}
