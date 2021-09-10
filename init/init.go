package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/habte/Project/MaidLink/internal/pkg/db"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/pkg"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var database *mongo.Database

// database varable gets the mongodb database instance.
// var database *mongo.Database
var once sync.Once

func init() {
	godotenv.Load("../conf.env")
	once.Do(
		func() {
			database = db.ConnectMongodb()
			if database == nil {
				os.Exit(1)
			}
		},
	)
}

func main() {
	database = db.ConnectMongodb()
	if database == nil {
		return
	}
	// collection := database.Collection("admin")
	// userCollection := database.Collection("user")
	pass, er := pkg.HashPassword(os.Getenv("MAID_LINK_DEFAULT_PASSWORD"))
	if er != nil {
		return
	}
	admin := &model.Admin{
		// ID: "61182e98ebeeebbee33314b9",
		User: &model.User{
			// ID:       "61182e98ebeeebbee33314b9",
			Username: "habte",
			Email:    "habteadnew@gmail.com",
			Password: pass,
			Role:     model.ADMIN,
			// ImageUrl: ,
			// CreatedBy: "61182e98ebeeebbee33314b9",
		},
		// CreatedBy: "61182e98ebeeebbee33314b9",
	}
	userInsertResult, era := database.Collection(model.SUSER).InsertOne(context.TODO(), admin.User)
	if era != nil {
		println("Error While Inserting User ", era.Error())
		return
	}
	// id := pkg.ObjectIDFromInsertResult(userInsertResult)
	admin.BsonID = userInsertResult.InsertedID.(primitive.ObjectID)
	admin.CreatedBy = pkg.RemoveObjectIDPrefix(pkg.ObjectIDFromInsertResult(userInsertResult))

	_, er = database.Collection(model.SADMIN).InsertOne(context.TODO(), admin)
	// if er != nil {
	// 	println("Error While Inserting the admin ... ", er.Error())
	// 	return
	// }

	// objectID := pkg.ObjectIDFromInsertResult(insertResult)
	// println(insertResult.InsertedID.(primitive.ObjectID).(string))

	// CREATING AN INDEX FOR A COLLECTION "MAIDS".
	index := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "carrers", Value: bsonx.String("text")},
				{Key: "address", Value: bsonx.String("text")},
				{Key: "bio", Value: bsonx.String("text")},
				{Key: "phone", Value: bsonx.String("text")},
				{Key: "works", Value: bsonx.String("text")},
				// {Key: "carrers", Value: bsonx.Array(bsonx.Arr{[]bsonx.Val{bsontype.Array.String(), byte("text"}})},
			},
		},
	}
	indexUser := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "username", Value: bsonx.String("text")},
				{Key: "email", Value: bsonx.String("text")},
				{Key: "role", Value: bsonx.Int32(1)},
				// {Key: "carrers", Value: bsonx.Array(bsonx.Arr{[]bsonx.Val{bsontype.Array.String(), byte("text"}})},
			},
		},
	}
	/*
		bson.D{{"myFirstTextField", "text"},{"mySecondTextField", "text"}}*/
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, errIndex := database.Collection(model.SMAID).Indexes().CreateMany(context.TODO(), index, opts)
	if errIndex != nil {
		println(errIndex.Error())
		panic(errIndex)
	}

	// -----------------------------------------------
	opts = options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, errIndex = database.Collection(model.SUSER).Indexes().CreateMany(context.TODO(), indexUser, opts)
	if errIndex != nil {
		println("Creating an index in the user : ", errIndex.Error())
		panic(errIndex)
	}
}
