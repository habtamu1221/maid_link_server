package admin

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IAdminService interface {
	CreateAdmin(context.Context) *model.Admin
	// GetAdmin  uses 'user_id' value in the context then, returns instance of admin model
	GetAdmin(context.Context) *model.Admin
	// GetMyAdmins  "admin_id" as a parameter
	GetMyAdmins(context.Context) []*model.Admin
	// DeleteMyAdmin "admin_id" <<-- Represents the Admin ID
	// "user_id" <<-- This represents teh owner admins id.
	DeleteMyAdmin(context.Context) bool
}
