package clientRepo

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/client"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClientRepository struct representing the client repository of mongo dataase
type ClientRepo struct {
	DB *mongo.Database //
}

// NewClientRepository returns a new client repository class performing on a mongo database
func NewClientRepository(db *mongo.Database) client.IClientRepository {
	return &ClientRepo{
		DB: db,
	}
}

func (clientr *ClientRepo) CreateClient(contex context.Context) (*model.Client, error) {
	client := contex.Value("client").(*model.Client)
	if oid, er := primitive.ObjectIDFromHex(client.User.ID); er == nil {
		client.BsonID = oid //
		client.ID = client.User.ID
		_, er := clientr.DB.Collection(model.SCLIENT).InsertOne(contex, client)
		return client, er
	}
	return nil, nil
}

func (clientr *ClientRepo) GetClient(conte context.Context) (*model.Client, error) {
	userID := conte.Value("user_id").(string)
	client := &model.Client{}
	if oid, er := primitive.ObjectIDFromHex(userID); er == nil {
		er = clientr.DB.Collection(model.SCLIENT).FindOne(conte, bson.D{{"_id", oid}}).Decode(client)
		client.ID = pkg.RemoveObjectIDPrefix(client.BsonID.String())
		return client, er
	} else {
		return nil, er
	}
}

// UpdateClient ...
func (clientr *ClientRepo) UpdateMyMaids(conte context.Context) ([]string, error) {
	mymaids := conte.Value("my_maids").([]string)
	if oid, er := primitive.ObjectIDFromHex(conte.Value("user_id").(string)); er == nil {
		if uc, er := clientr.DB.Collection(model.SCLIENT).UpdateOne(conte,
			bson.D{{"_id", oid}},
			bson.D{{"$set", bson.D{{"mymaids", mymaids}}}}); er == nil && uc.ModifiedCount > 0 {
			return mymaids, er
		} else {
			println(er.Error())
			return nil, er
		}
	} else {
		return nil, er
	}
}

// MyMaidsWhichIPayedFor "user_id" returns []*string
func (clientr *ClientRepo) MyMaidsWhichIPayedForString(conte context.Context) ([]string, error) {
	client := &model.Client{}
	if oid, er := primitive.ObjectIDFromHex(conte.Value("user_id").(string)); er == nil {
		if er = clientr.DB.Collection(model.SCLIENT).FindOne(conte, bson.D{{"_id", oid}}).Decode(client); er == nil {
			return client.MyMaids, er
		}
	}
	return nil, nil
}
