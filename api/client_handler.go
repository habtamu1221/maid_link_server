package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/habte/Project/MaidLink/internal/pkg/admin"
	"github.com/habte/Project/MaidLink/internal/pkg/client"
	"github.com/habte/Project/MaidLink/internal/pkg/maid"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/internal/pkg/session"
	"github.com/habte/Project/MaidLink/internal/pkg/user"
	"github.com/habte/Project/MaidLink/pkg"
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
			response.Write(pkg.GetJson(&model.ShortError{Err: "User account already exist"}))
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
			conte := context.WithValue(request.Context(), "maid_id", request.Context().Value("user_id").(string))
			client = clienth.MService.GetMaid(conte)
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
	response.Header().Set("Content-Type", "application/json")
	input := &model.PayForMaid{}
	jsonDec := json.NewDecoder(request.Body)
	if decerr := jsonDec.Decode(input); decerr == nil {
		// check the existance of the creadit care with this account number and
		// check whether the password is simmilar to the input one
		// check the existance of the maid with this ID
		// check whetehr the maid is in his my maids list
		// update my maids list of the client
		// return the maid full profile.
		if file, err := os.Open("dummyData/accounts.csv"); err == nil {
			defer file.Close()
			reader := csv.NewReader(file)
			writer := csv.NewWriter(file)
			if values, err := reader.ReadAll(); err == nil {
				for a, account := range values {
					if strings.Trim(input.AccNO, " ") == strings.Trim(account[0], " ") &&
						strings.Trim(input.Password, " ") == strings.Trim(account[1], " ") {
						if client := clienth.Service.GetClient(request.Context()); client != nil {
							if client.MyMaids == nil {
								client.MyMaids = []string{}
							}
							conte := request.Context()
							// check whether the aid exist or not
							userID := conte.Value("user_id").(string)
							conte = context.WithValue(conte, "maid_id", input.MaidID)
							if maid := clienth.MService.GetMaid(conte); maid != nil {
								exist := false //
								for _, id := range client.MyMaids {
									if id == input.MaidID {
										exist = true
									}
								}
								client.MyMaids = append(client.MyMaids, input.MaidID)
								conte = context.WithValue(conte, "user_id", userID)
								conte = context.WithValue(conte, "my_maids", client.MyMaids)
								if exist {
									response.Write(pkg.GetJson(map[string]interface{}{"maid_id": input.MaidID, "maid": maid}))
									return
								}
								if client := clienth.Service.UpdateMyMaids(conte); client != nil {
									values[a][2] = func() string {
										val, _ := strconv.Atoi(values[a][2])
										return strconv.Itoa((val - 1))
									}()
									println(values[a])
									if err = writer.WriteAll(values); err == nil {
										response.Write(pkg.GetJson(map[string]interface{}{"maid_id": input.MaidID, "maid": maid}))
										return
									} else {
										response.WriteHeader(http.StatusInternalServerError)
										response.Write(pkg.GetJson(&model.ShortError{"Internal Server Error "}))
										return
									}
								} else {
									response.WriteHeader(http.StatusInternalServerError)
									response.Write(pkg.GetJson(&model.ShortError{"Internal Server Error "}))
									return
								}
							} else {
								response.WriteHeader(http.StatusNotFound)
								response.Write(pkg.GetJson(&model.ShortError{"Maid not found"}))
								return
							}
						} else {
							response.WriteHeader(http.StatusUnauthorized)
							response.Write(pkg.GetJson(&model.ShortError{"UnAuthorized user"}))
							return
						}
					}
				}
				response.WriteHeader(http.StatusNotFound)
				response.Write(pkg.GetJson(&model.ShortError{"invalid account or password"}))
			}
		} else {
			response.Write(pkg.GetJson(&model.ShortError{"TRANSACTION: Internal Server Error "}))
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

// SearchForMaid  .... .
func (clienth *ClientHandler) SearchForMaid(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	println("CALLEDDDDD")
	offset, er := strconv.Atoi(request.FormValue("offset"))
	if er != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(request.FormValue("offset"))
	if err != nil || limit == offset {
		limit = offset + 2
	}
	if text := request.FormValue("q"); text != "" {
		println("Search q=", text)
		conte := context.WithValue(request.Context(), "q", text)
		conte = context.WithValue(conte, "offset", offset)
		conte = context.WithValue(conte, "limit", limit)
		if maids := clienth.MService.SearchMaids(conte); maids != nil {
			response.Write(pkg.GetJson(maids))
			return
		}
		response.WriteHeader(http.StatusNotFound)
		return
	}
	response.Write(pkg.GetJson(&model.ShortError{"invalid query"}))
	response.WriteHeader(http.StatusNotFound)
}
