package user

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IUserService interface {
	GetUser(context context.Context) *model.User
}
