package session

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
)

var jwtKey string
var cookiename string

type SessionHandler struct {
}

func NewSessionHandler() *SessionHandler {
	jwtKey = os.Getenv("JWT_KEY")
	cookiename = os.Getenv("COOKIE_NAME")
	return &SessionHandler{}
}

func (sess *SessionHandler) SaveSession(response http.ResponseWriter, request *http.Request, session *model.Session) (string, bool) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(24 * time.Hour)
	session.StandardClaims = jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
		// HttpOnly:  true,
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		println("-------------------HI-----------------------")
		return "", false
	}
	response.Header().Add("Authorization", "Bearer "+tokenString)
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	return tokenString, true
}

func (sess *SessionHandler) DeleteSession(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Authorization", "")
}

// GetSession returns the session object from token in the barier
func (sess *SessionHandler) GetSession(request *http.Request) *model.Session {
	token := request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil
	}
	session := &model.Session{}
	tkn, err := jwt.ParseWithClaims(token, session, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil
	}
	if !tkn.Valid {
		return nil
	}
	return session
}
