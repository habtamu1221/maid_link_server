package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/pkg"
)

type AdminHandler struct {
	SessionHandler *session.SessionHandler
	Service        admin.IAdminService
}

// NewAdminHandler for creating an admin ...
func NewAdminHandler(sessionh *session.SessionHandler, ser admin.IAdminService) *AdminHandler {
	return &AdminHandler{
		SessionHandler: sessionh,
		Service:        ser,
	}
}

// CreateAdmin admin handler method
// METHOD : POST
// AUTHORIZATION : Admins Only
func (adminh *AdminHandler) RegisterAdmin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	/*
		Expected JSON :
		{
			"username":
			"email":
			"password":
		}
	*/
	inpu := &struct {
		Username string
		Email    string
		Password string
	}{}
	decoder := json.NewDecoder(request.Body)
	if er := decoder.Decode(inpu); er == nil && inpu.Email != "" && inpu.Password != "" && inpu.Username != "" {
		user := &model.User{Username: inpu.Username, Email: inpu.Email, Password: inpu.Password}
		user.Role = model.ADMIN
		admin := &model.Admin{User: user, CreatedBy: request.Context().Value("session").(*model.Session).UserID}
		admin.CreatedBy = request.Context().Value("session").(*model.Session).UserID
		ncont := context.WithValue(request.Context(), "admin", admin)
		if admin = adminh.Service.CreateAdmin(ncont); admin != nil {
			response.WriteHeader(http.StatusCreated)
			response.Write(pkg.GetJson(&model.AdminInfo{Success: true, Admin: admin}))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(pkg.GetJson(&model.ShortError{Err: "Internal Server ERROR"}))
		return
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "bad request"}))
		return
	}
}
