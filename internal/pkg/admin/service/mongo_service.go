package admin_service

import "github.com/samuael/Project/MaidLink/internal/pkg/admin"

// AdminService struct representing.
type AdminService struct {
	Repo admin.IAdminRepository
}

// NewAdminService returns a new admin service instance
func NewAdminService(repo admin.IAdminRepository) admin.IAdminService {
	return &AdminService{
		Repo: repo,
	}
}
