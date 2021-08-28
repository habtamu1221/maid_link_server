package main

import (
	"context"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/samuael/Project/MaidLink/internal/pkg/db"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/pkg"
	"go.mongodb.org/mongo-driver/mongo"
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
	userCollection := database.Collection("user")
	pass, er := pkg.HashPassword(os.Getenv("MAID_LINK_DEFAULT_PASSWORD"))
	if er != nil {
		return
	}
	println("The Password is : " + pass)
	admin := &model.Admin{
		ID: "",
		User: &model.User{
			Username: "Samuael",
			Email:    "samuaeladnew@gmail.com",
			Password: pass,
			Role:     model.ADMIN,
			// ImageUrl: ,
		},
		CreatedBy: "",
	}
	userInsertResult, era := userCollection.InsertOne(context.TODO(), admin.User)
	if era != nil {
		println("Error While Inserting User ", era.Error())
		return
	}
	id := pkg.ObjectIDFromInsertResult(userInsertResult)
	// admin.UserID = id
	// insertResult, er := collection.InsertOne(context.TODO(), admin)
	// if er != nil {
	// 	println("Error While Inserting the admin ... ", er.Error())
	// 	return
	// }

	// objectID := pkg.ObjectIDFromInsertResult(insertResult)
	println(id)
}
