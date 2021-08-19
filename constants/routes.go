// Package constant ...
package constant

import "github.com/samuael/Project/MaidLink/internal/pkg/model"

//List of roues will be listed here...
var Routes = map[string]*model.Permission{
	"/admin": &model.Permission{
		Roles:   []string{model.SADMIN, model.SMAID}, //
		Methods: []string{"POST", "GET"},
	},
}
