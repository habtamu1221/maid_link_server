package api

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/samuael/Project/MaidLink/internal/pkg/maid"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/session"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"github.com/samuael/Project/MaidLink/pkg"
)

type MaidHandler struct {
	Session  *session.SessionHandler
	Service  maid.IMaidService
	UService user.IUserService
}

func NewMaidHandler(sess *session.SessionHandler, ser maid.IMaidService, uservice user.IUserService) *MaidHandler {
	return &MaidHandler{
		Session:  sess,
		Service:  ser,
		UService: uservice,
	}
}

// RegisterMaid .. create a maid
func (maidh *MaidHandler) Registermaid(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	maidInfo := &struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
		Email    string `json:"email"`
	}{}
	jsonDecoder := json.NewDecoder(request.Body)
	if decEr := jsonDecoder.Decode(maidInfo); decEr == nil && maidInfo.Username != "" && maidInfo.Phone != "" && maidInfo.Email != "" {
		contex := context.WithValue(request.Context(), "email", maidInfo.Email)
		if emailExists := maidh.UService.CheckEmailExistance(contex); !emailExists {
			user := &model.User{
				Username: maidInfo.Username,
				Email:    maidInfo.Email,
				Password: func() string {
					pwd, _ := pkg.HashPassword("1234")
					return pwd
				}(),
				Role: model.MAID,
			}
			maid := &model.Maid{
				ProfileImages: []string{},
				Phone:         maidInfo.Phone,
				User:          user,
				CreatedBy:     contex.Value("session").(*model.Session).UserID,
				Address:       maidInfo.Address,
			}
			contex = context.WithValue(contex, "maid", maid)
			if maid = maidh.Service.CreateMaid(contex); maid != nil {
				response.WriteHeader(http.StatusCreated)
				response.Write(pkg.GetJson(&model.UserInfo{Success: true, User: maid}))
				return
			}
			maid.User.Password = ""
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write(pkg.GetJson(&model.ShortError{Err: "Account with this email already exist!"}))
			return
		}
	}
	response.WriteHeader(http.StatusBadRequest)
	response.Write(pkg.GetJson(&model.ShortError{Err: "invalid input data"}))
}

// AddProfileImage handler fuction ...
func (maidh *MaidHandler) AddProfileImage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var header *multipart.FileHeader
	var erro error
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
		newName := "images/posts/" + pkg.GenerateRandomString(5, pkg.CHARACTERS) + "." + pkg.GetExtension(header.Filename)
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
		// oldImage = maidh.UService.GetImageUrl(request.Context())
		_, er := io.Copy(newImage, image)
		if er != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		ncon := context.WithValue(request.Context(), "user_id", request.Context().Value("session").(*model.Session).UserID)
		ncon = context.WithValue(ncon, "image_url", newName)
		_, success := maidh.Service.AddImageUrl(ncon)
		if success {
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

// RemoveProfileImage  a methdo to remove a profile
func (maidh *MaidHandler) RemoveProfileImage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	input := &struct {
		Imageurl string `json:"image_url"`
	}{}
	jsonDecode := json.NewDecoder(request.Body)
	if err := jsonDecode.Decode(input); err == nil {
		contx := context.WithValue(request.Context(), "image_url", input.Imageurl)
		if success := maidh.Service.RemoveProfileImage(contx); success {
			response.Write(pkg.GetJson(&model.Success{Success: true}))
			return
		} else {
			response.WriteHeader(http.StatusNotModified)
			// response.Write(pkg.GetJson(&model.Success{false}))
			return
		}
	} else {
		response.WriteHeader(http.StatusBadRequest)
	}
}

// GetProfileImages .. this is form value and Method GET
// user_id : the ID of the owner of the profile
func (maidh *MaidHandler) GetProfileImages(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if profileImage := request.FormValue("user_id"); profileImage != "" {
		println(profileImage)
		contex := context.WithValue(request.Context(), "user_id", profileImage)
		if images := maidh.Service.GetImageUrls(contex); images != nil {
			response.Write(pkg.GetJson(images))
			return
		}
		response.WriteHeader(http.StatusNotFound)
		return
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "invalid request"}))
		return
	}
}

// GetUser a method :
func (maidh *MaidHandler) GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if profileImage := request.FormValue("user_id"); profileImage != "" {
		println(profileImage)
		contex := context.WithValue(request.Context(), "user_id", profileImage)
		if user := maidh.Service.GetUser(contex); user != nil {
			response.Write(pkg.GetJson(user))
			return
		}
		response.WriteHeader(http.StatusNotFound)
		return
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{Err: "invalid request"}))
		return
	}
}

func (maidh *MaidHandler) CreateWork(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	input := &model.Work{}
	jsonDecode := json.NewDecoder(request.Body)
	if err := jsonDecode.Decode(input); err == nil {
		conte := request.Context()
		session := conte.Value("session").(*model.Session)
		conte = context.WithValue(conte, "user_id", session.UserID)
		conte = context.WithValue(conte, "work", input)
		if work := maidh.Service.CreateWork(conte); work != nil {
			response.WriteHeader(http.StatusCreated)
			response.Write(pkg.GetJson(work))
			return
		}
		response.WriteHeader(http.StatusNotModified)
		response.Write(pkg.GetJson(&model.ShortError{Err: "Not Modified"}))
		return
	}
	response.WriteHeader(http.StatusBadRequest)
}

// DeleteWork  DELETE AUTHORIZED only for MAIDS , input  no=2  as a parameter
func (maidh *MaidHandler) DeleteWork(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if no, err := strconv.Atoi(request.FormValue("no")); err == nil && no > 0 {
		session := request.Context().Value("session").(*model.Session)
		conte := context.WithValue(request.Context(), "user_id", session.UserID)
		conte = context.WithValue(conte, "work_no", no)
		if success := maidh.Service.DeleteWork(conte); success {
			response.Write(pkg.GetJson(&model.ShortSuccess{"succesfuly deleted work number " + strconv.Itoa(no)}))
			return
		}
		response.WriteHeader(http.StatusNotFound)
		response.Write(pkg.GetJson(&model.Success{false}))
		return
	}
	response.Write(pkg.GetJson(&model.ShortError{"invalid input"}))
	response.WriteHeader(http.StatusBadRequest)
}

// UpdateWork method to update the work using the work ID
func (maidh *MaidHandler) UpdateWork(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	input := &model.Work{}
	jsonDecode := json.NewDecoder(request.Body)
	if err := jsonDecode.Decode(input); err == nil {
		conte := context.WithValue(request.Context(), "work", input)
		if input = maidh.Service.UpdateWork(conte); input != nil {
			response.Write(pkg.GetJson(input))
			return
		}
		response.WriteHeader(http.StatusNotFound)
		response.Write(pkg.GetJson(&model.ShortError{"update was not succesful!"}))
		return
	}
	response.WriteHeader(http.StatusBadRequest)
	response.Write(pkg.GetJson(&model.ShortError{"bad input"}))
}

// GetWorks ...
func (maidh *MaidHandler) GetWorks(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID := request.FormValue("id")
	conte := request.Context()
	if userID != "" {
		conte = context.WithValue(conte, "user_id", userID)
	}
	works := maidh.Service.GetWorks(conte)
	if works != nil {
		response.Write(pkg.GetJson(works))
		return
	}
	response.Write(pkg.GetJson([]*model.Work{}))
	response.WriteHeader(http.StatusNotFound)
}

func (maidh *MaidHandler) UpdateMaid(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

}

// GetAdminMaids returns the list of maids created by this admin...
func (maidh *MaidHandler) GetAdminMaids(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	conte := request.Context()
	adminID := request.FormValue("admin_id")
	if adminID != "" {
		conte = context.WithValue(conte, "user_id", adminID)
	}
	maids := maidh.Service.GetMyMaids(conte)
	if maids == nil {
		response.Write(pkg.GetJson([]interface{}{}))
		return
	}
	response.Write(pkg.GetJson(maids))
}

// MaidUpdateProfile handler function to update the maids profile
// this update will be applicabel only for the maid document the other additional Updates are to be
// performed using user update request.
func (maidh *MaidHandler) MaidUpdateProfile(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	maidu := model.MaidUpdate{}
	jsonDec := json.NewDecoder(request.Body)

	if decErro := jsonDec.Decode(maidu); decErro == nil || maidu.ID == "" {
		// Get The Maid Using the ID of the maid
		conte := request.Context()
		if maid := maidh.Service.GetMaid(conte); maid != nil {
			maid.Phone = func() string {
				if maidu.Phone != "" || len(maidu.Phone) >= 10 || func() bool {
					if _, err := strconv.Atoi(strings.TrimPrefix("+", maidu.Phone)); err == nil {
						return true
					} else {
						return false
					}
				}() {
					return maidu.Phone
				} else {
					return maid.Phone
				}
			}()
			maid.Address = func() string {
				if maidu.Address != "" {
					return maidu.Address
				} else {
					return maid.Address
				}
			}()
			maid.Bio = func() string {
				if maidu.Bio != "" {
					return maidu.Bio
				} else {
					return maid.Bio
				}
			}()
			conte = context.WithValue(conte, "maid", maid)
			if maid = maidh.Service.UpdateMaid(conte); maid != nil {
				// if the maid is not null then do this ...
				maidu.ID = maid.ID
				maidu.Phone = maid.Phone
				maidu.Address = maid.Address
				maidu.Bio = maid.Bio
				response.Write(pkg.GetJson(maidu))
				return
			} else {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write(pkg.GetJson(&model.ShortError{"Internal server error "}))
				return
			}
		} else {
			response.WriteHeader(http.StatusNotFound)
			response.Write(pkg.GetJson(&model.ShortError{"Resource not found.."}))
			return
		}
		// maid.ID = maidu.ID

		// maid.Username = func()string{ if maidu.Username != "" {return maidu.Username}return maid.Username  }()
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(pkg.GetJson(&model.ShortError{"Invalid json input!"}))
		return
	}
	// decoding a json
}
