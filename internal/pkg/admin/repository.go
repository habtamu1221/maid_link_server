package admin

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/model"
)

type IAdminRepository interface {
	CreateAdmin(context.Context) (*model.Admin, error)
	GetAdmin(context.Context) (*model.Admin, error)
	GetMyAdmins(context.Context) ([]*model.Admin, error)
	DeleteMyAdmin(context.Context) error
}
