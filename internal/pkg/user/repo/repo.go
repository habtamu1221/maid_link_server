package user_repo

import (
	"context"
	"errors"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Database
}

func NewUserRepo(db *mongo.Database) user.IUserRepo {
	return &UserRepo{DB: db}
}

func (userrepo UserRepo) GetUserByEmail(context context.Context) (*model.User, error) {
	email := context.Value("email").(string)
	user := &model.User{}
	er := userrepo.DB.Collection(model.SUSER).FindOne(context, bson.D{{"email", email}}).Decode(user)
	if er != nil {
		return nil, er
	}
	return user, nil
}

func (userrepo UserRepo) GetUserByID(context context.Context) (*model.User, error) {
	userID := context.Value("user_id").(string)
	oid, er := primitive.ObjectIDFromHex(userID)
	if er != nil {
		println(er.Error())
		return nil, er
	}
	user := &model.User{}
	er = userrepo.DB.Collection(model.SUSER).FindOne(context, bson.D{{"_id", oid}}).Decode(user)
	if er != nil {
		return nil, er
	}
	return user, nil
}

func (userrepo *UserRepo) UpdateUser(context context.Context) (*model.User, error) {
	user := context.Value("user").(*model.User)
	if oid, er := primitive.ObjectIDFromHex(user.ID); er == nil {
		userrepo.DB.Collection(model.SUSER).UpdateOne(context, bson.D{{"_id", oid}}, bson.D{{"$set",
			bson.D{
				{"username", user.Username},
				{"password", user.Password},
				{"email", user.Email},
				{"imageurl", user.ImageUrl},
				{"role", user.Role},
			}}})
		return user, nil
	} else {
		return nil, er
	}
}

func (userrepo *UserRepo) ChangePassword(context context.Context) (*model.User, error) {
	user := context.Value("user").(*model.User)
	if oid, er := primitive.ObjectIDFromHex(user.ID); er == nil {
		userrepo.DB.Collection(model.SUSER).UpdateOne(context, bson.D{{"_id", oid}}, bson.D{{"$set",
			bson.D{
				{"password", user.Password},
			}}})
		return user, nil
	} else {
		return nil, er
	}
}
func (userrepo *UserRepo) ChangeImageUrl(context context.Context) error {
	userID := context.Value("user_id").(string)
	imageUrl := context.Value("image_url").(string)
	if oid, er := primitive.ObjectIDFromHex(userID); er == nil {
		userrepo.DB.Collection(model.SUSER).UpdateOne(context,
			bson.D{{"_id", oid}},
			bson.D{{"$set",
				bson.D{
					{"image_url", imageUrl},
				}}})
		return nil
	} else {
		return er
	}
}

func (userrepo *UserRepo) GetImageUrl(context context.Context) (string, error) {
	userid := context.Value("session").(*model.Session).UserID
	user := &model.User{}
	if oid, er := primitive.ObjectIDFromHex(userid); er == nil {
		er := userrepo.DB.Collection(model.SUSER).FindOne(context,
			bson.D{{"_id", oid}}).Decode(user)
		if er == nil {
			return user.ImageUrl, nil
		}
	}
	return "", errors.New("Internal Error ")
}
