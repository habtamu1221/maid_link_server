package admin_repository

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/pkg"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminRepository struct representing the admin   repository for mongodb
type AdminRepo struct {
	DB *mongo.Database //
}

func NewAdminRepository(db *mongo.Database) admin.IAdminRepository {
	return &AdminRepo{
		DB: db,
	}
}

// CreateAdmin a method which merely create an admin instance using ony the created By and user id Information
// the rest of admin information is to be handled by by the create user function
// the user information is not populated in this instance.
func (adminr *AdminRepo) CreateAdmin(context context.Context) (*model.Admin, error) {
	admin := context.Value("admin").(*model.Admin)
	nid := pkg.ObjectIDFromString(admin.User.ID)
	admin.User.ID = nid
	if isertID, er := adminr.DB.Collection(model.SADMIN).InsertOne(context, admin); er == nil {
		admin.ID = pkg.ObjectIDFromInsertResult(isertID)
		return admin, nil
	} else {
		return nil, er
	}
}
