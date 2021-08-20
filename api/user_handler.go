package api

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"github.com/samuael/Project/MaidLink/pkg"
)

// UserHandler ...
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
	user := userhandler.Service.GetUserByEmail(context)
	if user == nil {
		response.WriteHeader(401)
		response.Write(pkg.GetJson(&model.Error{Message: "error", Reason: "User Not Authorized"}))
		return
	}
	// ------- check the password
	if !(pkg.CompareHash(user.Password, loginData.Password)) {
		response.WriteHeader(http.StatusNotFound)
		response.Write(pkg.GetJson(&model.ShortError{Err: "Invalid Username or Password"}))
		return
	}
	user.Password = ""
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
	response.Header().Set("Content-Type", "application/json")
	val := &struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm"`
	}{}
	jsonDecode := json.NewDecoder(request.Body)
	if err := jsonDecode.Decode(val); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "Invalid Input"}))
		return
	} else if val.Password != val.ConfirmPassword {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "Password Confirmation Error"}))
		return
	}

	if session := userhandler.SessionHandler.GetSession(request); session != nil {
		ncont := context.WithValue(request.Context(), "user_id", session.UserID)
		if user := userhandler.Service.GetUserByID(ncont); user != nil {
			if pkg.CompareHash(user.Password, val.Password) {
				response.WriteHeader(http.StatusNotAcceptable)
				response.Write(pkg.GetJson(&model.ShortSuccess{Msg: "No Change is Made"}))
				return
			}
			if newPassHash, err := pkg.HashPassword(val.ConfirmPassword); err == nil && newPassHash != "" {
				user.Password = newPassHash
				// save the user to Database
				context := context.WithValue(request.Context(), "user", user)
				user = userhandler.Service.ChangePassword(context)
				if user == nil {
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				response.Write(pkg.GetJson(&model.PasswordChangeSuccess{NewPassword: val.ConfirmPassword}))
				return
			} else {

				response.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			response.WriteHeader(http.StatusNotFound)
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
	}
}

// ChangeProfilePicture  for user to change profile Picture ....
func (userhandler *UserHandler) ChangeProfilePicture(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var header *multipart.FileHeader
	var erro error
	var oldImage string
	erro = request.ParseMultipartForm(99999999999)
	if erro != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	// -----
	image, header, erro := request.FormFile("image")
	if erro != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()
	if pkg.IsImage(header.Filename) {
		newName := "images/profile/" + pkg.GenerateRandomString(5, pkg.CHARACTERS) + "." + pkg.GetExtension(header.Filename)
		var newImage *os.File
		if strings.HasSuffix(os.Getenv("ASSETS_DIRECTORY"), "/") {
			newImage, erro = os.Create(os.Getenv("ASSETS_DIRECTORY") + newName)
		} else {
			newImage, erro = os.Create(os.Getenv("ASSETS_DIRECTORY") + "/" + newName)
		}
		if erro != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer newImage.Close()
		oldImage = userhandler.Service.GetImageUrl(request.Context())
		_, er := io.Copy(newImage, image)
		if er != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		ncon := context.WithValue(request.Context(), "user_id", request.Context().Value("session").(*model.Session).UserID)
		ncon = context.WithValue(ncon, "image_url", newName)
		success := userhandler.Service.ChangeImageUrl(ncon)
		if success {
			if oldImage != "" {
				if strings.HasSuffix(os.Getenv("ASSETS_DIRECTORY"), "/") {
					er = os.Remove(os.Getenv("ASSETS_DIRECTORY") + oldImage)
				} else {
					er = os.Remove(os.Getenv("ASSETS_DIRECTORY") + "/" + oldImage)
				}
			}
			response.WriteHeader(http.StatusOK)
			response.Write(pkg.GetJson(&model.ShortSuccess{Msg: newName}))
			return
		}
		if strings.HasSuffix(os.Getenv("ASSETS_DIRECTORY"), "/") {
			er = os.Remove(os.Getenv("ASSETS_DIRECTORY") + newName)
		} else {
			er = os.Remove(os.Getenv("ASSETS_DIRECTORY") + "/" + newName)
		}
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.WriteHeader(http.StatusUnsupportedMediaType)
	}
}

// DeleteProfilePicture ...
func (userhandler *UserHandler) DeleteProfilePicture(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	imageUrl := userhandler.Service.GetImageUrl(request.Context())
	success := userhandler.Service.DeleteProfilePicture(request.Context())
	if success {
		if strings.HasSuffix(os.Getenv("ASSETS_DIRECTORY"), "/") {
			os.Remove(os.Getenv("ASSETS_DIRECTORY") + imageUrl)
		} else {
			os.Remove(os.Getenv("ASSETS_DIRECTORY") + "/" + imageUrl)
		}
		response.Write(pkg.GetJson(&model.ShortSuccess{Msg: "Succesfully Deleted"}))
		return
	} else {
		response.WriteHeader(http.StatusInternalServerError)
	}
}

// RegisterUser ...
func (userhandler *UserHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
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
		user.Role = model.CLIENT
		ncont := context.WithValue(request.Context(), "user", user)
		if user = userhandler.Service.CreateUser(ncont); user != nil {
			response.WriteHeader(http.StatusCreated)
			response.Write(pkg.GetJson(&model.UserInfo{Success: true, User: user}))
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
