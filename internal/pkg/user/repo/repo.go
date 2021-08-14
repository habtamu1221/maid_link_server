package user_repo

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Database
}

func NewUserRepo(db *mongo.Database) user.IUserRepo {
	return &UserRepo{DB: db}
}

func (userrepo UserRepo) GetUser(context context.Context) (*model.User, error) {
	email := context.Value("email").(string)
	user := &model.User{}
	er := userrepo.DB.Collection(model.SUSER).FindOne(context, bson.D{{"email", email}}).Decode(user)
	if er != nil {
		return nil, er
	}
	return user, nil
}
