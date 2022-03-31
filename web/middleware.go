package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
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

func (s Server) ValidBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getClaims(r)
		if err != nil {
			sendAuthErrorfJSON(w, "Invalid access token. Please login to continue")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s Server) GetUserIDTokenCtx(r *http.Request) string {
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

func (s Server) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := s.GetUserIDTokenCtx(r)
		if uid == "" {
			sendAuthErrorfJSON(w, "Invalid access token. Please login to continue")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s Server) ValidAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("API_AUTH_KEY")
		if key != os.Getenv("ACCESS_SECRET") {
			sendAuthErrorfJSON(w, "Invalid access token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

var busyUsers = map[string]bool{}

func (s Server) NoReentry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := s.GetUserIDTokenCtx(r)
		if uid == "" {
			sendAuthErrorfJSON(w, "Invalid access token. Please login to continue")
			return
		}

		if busyUsers[uid] {
			SendErrorfJSON(w, "you have a pending request, please try again later")
			return
		}
		busyUsers[uid] = true
		defer delete(busyUsers, uid)

		next.ServeHTTP(w, r)
	})
}

// ContractAddressCtx returns a http.HandlerFunc that embeds the value at the url
// part {contractAddress} into the request context.
func ContractAddressCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CtxContractAddress,
			chi.URLParam(r, "contractAddress"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetContractAddressCtx(r *http.Request) string {
	chartAxisType, ok := r.Context().Value(CtxContractAddress).(string)
	if !ok {
		log.Trace("chart axis type not set")
		return ""
	}
	return chartAxisType
}
