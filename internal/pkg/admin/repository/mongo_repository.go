package admin_repository

import (
	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminRepository struct representing the admin   repository for mongodb
type AdminRepository struct {
	DB *mongo.Database //
}

func NewAdminRepository(db *mongo.Database) admin.IAdminRepository {
	return &AdminRepository{
		DB: db,
	}
}
