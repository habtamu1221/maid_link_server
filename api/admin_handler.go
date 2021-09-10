package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/habte/Project/MaidLink/internal/pkg/admin"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/internal/pkg/session"
	"github.com/habte/Project/MaidLink/internal/pkg/user"
	"github.com/habte/Project/MaidLink/pkg"
)

type AdminHandler struct {
	SessionHandler *session.SessionHandler
	Service        admin.IAdminService
	UService       user.IUserService
}

// NewAdminHandler for creating an admin ...
func NewAdminHandler(sessionh *session.SessionHandler, ser admin.IAdminService, userser user.IUserService) *AdminHandler {
	return &AdminHandler{
		SessionHandler: sessionh,
		Service:        ser,
		UService:       userser,
	}
}

/*
	RegisterAdmins :: ADMINS   :: POST
*/

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

		conte := context.WithValue(request.Context(), "email", inpu.Email)
		if user := adminh.UService.GetUserByEmail(conte); user != nil {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write(pkg.GetJson(&model.ShortError{"User already exist!"}))
			return
		}
		pass, er := pkg.HashPassword(inpu.Password)
		if er != nil {
			response.Write(pkg.GetJson(&model.ShortError{"invalid password please try other password"}))
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		user := &model.User{Username: inpu.Username, Email: inpu.Email, Password: pass}
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

func (adminh *AdminHandler) GetMyAdmins(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	conte := context.WithValue(request.Context(), "admin_id", request.Context().Value("user_id").(string))
	if admins := adminh.Service.GetMyAdmins(conte); admins != nil {
		response.Write(pkg.GetJson(admins))
		return
	} else {
		response.WriteHeader(http.StatusNotFound)
		response.Write(pkg.GetJson(&model.ShortError{"no admins found "}))
		return
	}
}

func (adminh *AdminHandler) DeleteAdmin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if adminid := request.FormValue("admin_id"); adminid != "" {
		conte := request.Context()
		conte = context.WithValue(conte, "admin_id", adminid)
		if success := adminh.Service.DeleteMyAdmin(conte); success {
			response.Write(pkg.GetJson(&model.ShortSuccess{"deleted succesfuly"}))
			return
		} else {
			response.Write(pkg.GetJson(&model.ShortError{"Can't Delete the admin instance"}))
			response.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{"Invallid admin id "}))
	}
}
