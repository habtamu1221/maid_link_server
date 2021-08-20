package admin_service

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
)

// AdminService struct representing.
type AdminService struct {
	Repo  admin.IAdminRepository
	URepo user.IUserRepo
}

// NewAdminService returns a new admin service instance
func NewAdminService(repo admin.IAdminRepository, urepo user.IUserRepo) admin.IAdminService {
	return &AdminService{
		Repo:  repo,
		URepo: urepo,
	}
}

func (adminser *AdminService) CreateAdmin(contex context.Context) *model.Admin {
	// first of all create a user
	admin := contex.Value("admin").(*model.Admin)
	contex = context.WithValue(contex, "user", admin.User)
	if user, er := adminser.URepo.CreateUser(contex); er == nil {
		admin.User = user
		admin.ID = user.ID
		admin, er := adminser.Repo.CreateAdmin(contex)
		if er == nil {
			return admin
		}
		contex = context.WithValue(contex, "user_id", user.ID)
		adminser.URepo.RemoveUser(contex)
		return nil
	}
	return nil
}
