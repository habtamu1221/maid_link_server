// Package constant ...
package constant

import (
	"net/http"
	"strings"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

//List of roues will be listed here...
var Routes = map[string]*model.Permission{
	"/admin/new/": &model.Permission{
		Roles:   []int{model.ADMIN}, //
		Methods: []string{"POST"},
	},
	"/maid/new/": &model.Permission{
		Roles: []int{
			model.ADMIN,
		},
		Methods: []string{"POST"},
	},
	"/maid/profile/add/": &model.Permission{
		Roles: []int{
			model.MAID,
		},
		Methods: []string{"POST"},
	},
	"/maid/profile/remove/": &model.Permission{
		Roles: []int{
			model.MAID,
		},
		Methods: []string{"DELETE"},
	},
	"/maid/work/new/": &model.Permission{
		Roles: []int{
			model.MAID,
		},
		Methods: []string{"POST"},
	},
	"/maid/work/": &model.Permission{
		Roles: []int{
			model.MAID,
		},
		Methods: []string{"DELETE", "PUT"},
	},
	"/admin/maids/": &model.Permission{
		Roles: []int{
			model.ADMIN,
		},
		Methods: []string{"GET"},
	},

	// "/client/profile/": &model.Permission{
	// 	Roles: []int{
	// 		model.CLIENT,
	// 	},
	// 	Methods: []string{"GET"},
	// },
	"/maid/info/payment/": &model.Permission{
		Roles: []int{
			model.CLIENT,
		},
		Methods: []string{"GET", "POST"},
	},
	"/maid/rate/": &model.Permission{
		Roles: []int{
			model.CLIENT,
		},
		Methods: []string{"GET", "POST"},
	},
	"/admins/": &model.Permission{
		Roles: []int{
			model.ADMIN,
		},
		Methods: []string{"GET", "POST"},
	},
	"/admin/": &model.Permission{
		Roles: []int{
			model.ADMIN,
		},
		Methods: []string{"DELETE", "POST"},
	},
}

// IsAuthorized ...
func IsAuthorized(request *http.Request) (bool, int) {
	path := request.URL.Path
	session := request.Context().Value("session").(*model.Session)
	var permission *model.Permission
	path = strings.TrimPrefix(path, "/api")
	if strings.HasSuffix(path, "/") {
		permission = Routes[path]
	} else {
		permission = Routes[path+"/"]
	}
	if permission == nil {
		return false, http.StatusNotFound
	}
	for _, rl := range permission.Roles {
		if session.Role == rl {
			for _, method := range permission.Methods {
				if request.Method == method {
					return true, http.StatusOK
				}
			}
			return false, http.StatusMethodNotAllowed
		}
	}
	return false, http.StatusUnauthorized
}
