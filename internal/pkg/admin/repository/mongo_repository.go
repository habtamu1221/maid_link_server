package admin_repository

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/admin"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if oid, err := primitive.ObjectIDFromHex(admin.User.ID); err == nil {
		admin.ID = admin.User.ID
		admin.BsonID = oid
		if isertID, er := adminr.DB.Collection(model.SADMIN).InsertOne(context, admin); er == nil {
			admin.ID = pkg.ObjectIDFromInsertResult(isertID)
			return admin, nil
		} else {
			return nil, er
		}
	} else {
		return nil, err
	}
}
func (adminr *AdminRepo) GetAdmin(conte context.Context) (*model.Admin, error) {
	userID := conte.Value("user_id").(string)
	admin := &model.Admin{}
	if oid, er := primitive.ObjectIDFromHex(userID); er == nil {
		er = adminr.DB.Collection(model.SADMIN).FindOne(conte, bson.D{{"_id", oid}}).Decode(admin)
		admin.ID = pkg.RemoveObjectIDPrefix(admin.BsonID.String())
		return admin, er
	} else {
		println(er.Error())
		return nil, er
	}
}
