package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"github.com/samuael/Project/MaidLink/pkg"
)

type UserHandler struct {
	Service        user.IUserService //
	SessionHandler *session.SessionHandler
}

// UserHandler
func NewUserHandler(service user.IUserService, sessHandler *session.SessionHandler) *UserHandler {
	return &UserHandler{
		Service:        service,
		SessionHandler: sessHandler,
	}
}

// UserLogin a a login end point for login
func (userhandler *UserHandler) UserLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	jsonDecoder := json.NewDecoder(request.Body)
	loginData := &model.LoginInput{}

	eror := jsonDecoder.Decode(loginData)
	if eror != nil {
		response.WriteHeader(401)
		response.Write(pkg.GetJson(&model.Error{Message: "error", Reason: "Invalid Input "}))
		return
	}
	context := context.WithValue(request.Context(), "email", loginData.Email)
	user := userhandler.Service.GetUser(context)
	if user == nil {
		response.WriteHeader(401)
		response.Write(pkg.GetJson(&model.Error{Message: "error", Reason: "User Not Authorized"}))
		return
	}
	sessionValue := &model.Session{
		UserID: user.ID,
		Role:   int(user.Role),
	}
	tokenString, success := userhandler.SessionHandler.SaveSession(response, request, sessionValue)
	if !success {
		response.WriteHeader(401)
		response.Write(pkg.GetJson(&model.Error{Message: "err", Reason: "unauthorized"}))
		return
	}
	successMessage := model.LoginSuccess{
		Token: tokenString,
		User:  user,
	}
	response.Write(pkg.GetJson(successMessage))
}

// ChangePassword  ...
func (userhandler *UserHandler) ChangePassword(response http.ResponseWriter, request *http.Request) {

}
