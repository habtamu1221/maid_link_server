package api

import (
	"context"
	"net/http"
	"time"

	constant "github.com/samuael/Project/MaidLink/constants"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/pkg"
)

type Auth struct {
	// session handler instance pointer
	SessionHandler *session.SessionHandler
}

// https://www.youtube.com/watch?v=UO98lJQ3QGI

func NewAuth(sess *session.SessionHandler) *Auth {
	return &Auth{sess}
}

func (auth *Auth) Authorize(handler http.Handler) http.Handler {
	// Authorize is a function which check whether the function request with a session is valid or not
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if permission := constant.Routes[request.URL.Path]; permission != nil {
			handler.ServeHTTP(response, request)
		}
		http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

// SetContext ...
func (auth *Auth) SetContext(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(respnse http.ResponseWriter, request *http.Request) {
		contex, _ := context.WithDeadline(request.Context(), time.Now().Add(time.Millisecond*500))
		request = request.WithContext(contex)
		handler.ServeHTTP(respnse, request)
	})
}

// Logout function api Logging out
// METHOD GET
// VAriables NONE
func (auth *Auth) Logout(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Authorization", "")
	response.WriteHeader(http.StatusOK)
	response.Write(pkg.GetJson(model.Success{Success: true}))
}

// LoggedIn checks whether the user is Authenticated or not
func (auth *Auth) IsLoggedInUser(request *http.Request) bool {
	session := auth.SessionHandler.GetSession(request)
	return session != nil
}

// Authenticated checks if a user has proper authority to access a give route
func (auth *Auth) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsLoggedInUser(r) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		session := auth.SessionHandler.GetSession(r)
		if session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		con := context.WithValue(r.Context(), "session", session)
		r = r.WithContext(con)
		next.ServeHTTP(w, r)
	})
}
