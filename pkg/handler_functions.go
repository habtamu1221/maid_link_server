package pkg

import (
	"encoding/json"

	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetJson(object interface{}) []byte {
	val, _ := json.Marshal(object)
	return val
}

// ObjectIDFromInsertResult returning an insert ID as a string
// // taking an input as *monogo.InsertOneResult if its not valid
// // it will return an empty string
func ObjectIDFromInsertResult(sires *mongo.InsertOneResult) string {
	if sires == nil {
		return ""
	}
	slices := RemoveObjectIDPrefix(sires.InsertedID.(primitive.ObjectID).String())
	return slices
}

func ObjectIDFromString(objid string) string {
	return "ObjectID(\"" + objid + "\")"
}

// RemoveObjectIDPrefix function returning the real internal Object ID from
// ObjectID prefixed object ID Result
// Example Input : ObjectID("5fe1b21d88b1deda65a9a507") :
// 		   OutPut : "5fe1b21d88b1deda65a9a507"
func RemoveObjectIDPrefix(objectid string) string {
	objectid = strings.TrimSuffix(strings.TrimPrefix(objectid, "ObjectID(\""), "\")")
	return objectid
}
