package api

import (
	"context"
	"net/http"
	"time"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/pkg"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) Authenticate(handler http.Handler) *http.Handler {
	return nil
}

func (auth *Auth) Authorize(handler http.Handler) http.Handler {
	return nil
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

// Login api Login function
// METHOD POST
// INPUT JSON    {
// email :
// password :
// }
