package admin

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IAdminRepository interface {
	CreateAdmin(context.Context) (*model.Admin, error)
}
