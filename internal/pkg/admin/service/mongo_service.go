package admin_service

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/admin"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/internal/pkg/user"
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
		contex = context.WithValue(contex, "admin", admin)
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

func (adminser *AdminService) GetAdmin(conte context.Context) *model.Admin {

	if admin, er := adminser.Repo.GetAdmin(conte); er == nil {
		// conte = context.WithValue(conte , "user_id"  , )
		if admin.User, er = adminser.URepo.GetUserByID(conte); er == nil && admin.User != nil {
			return admin
		}
	}
	return nil
}

func (adminser *AdminService) GetMyAdmins(conte context.Context) []*model.Admin {
	if admins, err := adminser.Repo.GetMyAdmins(conte); err == nil {
		for a, admin := range admins {
			conte = context.WithValue(conte, "user_id", admin.ID)
			println(admin.ID)
			admin.User, _ = adminser.URepo.GetUserByID(conte)
			admins[a] = admin
		}
		return admins
	}
	return nil
}

func (adminser *AdminService) DeleteMyAdmin(conte context.Context) bool {
	return adminser.Repo.DeleteMyAdmin(conte) == nil
}
