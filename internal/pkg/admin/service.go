package admin

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IAdminService interface {
	CreateAdmin(context.Context) *model.Admin
}
