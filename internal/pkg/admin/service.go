package admin

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IAdminService interface {
	CreateAdmin(context.Context) *model.Admin
	// GetAdmin  uses 'user_id' value in the context then, returns instance of admin model
	GetAdmin(context.Context) *model.Admin
}
