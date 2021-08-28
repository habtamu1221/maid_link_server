package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"github.com/samuael/Project/MaidLink/internal/pkg/client"
	"github.com/samuael/Project/MaidLink/internal/pkg/maid"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"github.com/samuael/Project/MaidLink/pkg"
)

type ClientHandler struct {
	Session  *session.SessionHandler
	Service  client.IClientService
	UService user.IUserService
	AService admin.IAdminService
	MService maid.IMaidService
}

func NewClientHandler(
	session *session.SessionHandler,
	ser client.IClientService,
	user user.IUserService,
	aser admin.IAdminService,
	mser maid.IMaidService) *ClientHandler {
	return &ClientHandler{
		Session:  session,
		Service:  ser,
		UService: user,
		AService: aser,
		MService: mser,
	}
}

// RegisterClient ...
func (clienth *ClientHandler) RegisterClient(response http.ResponseWriter, request *http.Request) {
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
		ncont := context.WithValue(request.Context(), "email", inpu.Email)
		if clienth.UService.CheckEmailExistance(ncont) {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write(pkg.GetJson(&model.ShortError{Err: "User Already Exist ... "}))
			return
		}
		passhash, er := pkg.HashPassword(inpu.Password)
		if er != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(pkg.GetJson(&model.ShortError{Err: "Please Enter another Password "}))
			return
		}
		user := &model.User{Username: inpu.Username, Email: inpu.Email, Password: passhash}
		user.Role = model.CLIENT
		client := &model.Client{User: user, MyMaids: []string{}}
		ncont = context.WithValue(request.Context(), "client", client)
		if client = clienth.Service.CreateClient(ncont); client != nil {
			response.WriteHeader(http.StatusCreated)
			response.Write(pkg.GetJson(&model.UserInfo{Success: true, User: client}))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(pkg.GetJson(&model.ShortError{Err: "registraction was not succesful!\nTry Again"}))
		return
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "bad request"}))
		return
	}
}

func (clienth *ClientHandler) GetClient(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	obj := &model.Profile{Success: false}
	session := request.Context().Value("session").(*model.Session)
	var era error
	var client interface{}
	switch session.Role {
	case model.ADMIN:
		{
			client = clienth.AService.GetAdmin(request.Context())
		}
	case model.CLIENT:
		{
			client = clienth.Service.GetClient(request.Context())
		}
	case model.MAID:
		{
			client = clienth.MService.GetMaid(request.Context())
		}
	}
	if era == nil && client != nil {
		obj.Body = client
		obj.Role = session.Role
		obj.Success = true
		response.Write(pkg.GetJson(obj))
		return
	}
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(pkg.GetJson(obj))
}

func (clienth *ClientHandler) PayForMaidInfo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Contexnt-Type", "application/json")
	input := &model.PayForMaid{}
	jsonDec := json.NewDecoder(request.Body)
	if decerr := jsonDec.Decode(input); decerr == nil {
		// check the existance of the creadit care with this account number and
		// check whether the password is simmilar to the input one
		// check the existance of the maid with this ID
		// check whetehr the maid is in his my maids list
		// update my maids list of the client
		// return the maid full profile.

	} else {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}
