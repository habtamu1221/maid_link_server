package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents all the users in the system
type User struct {
	ID       string `json:"id" bson:"_id,omitempty" `
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImageUrl string `json:"image_url" bson:"image_url" `
	Role     int8   `json:"role"`
}

// Client represents the client
type Client struct {
	ID      string             `json:"id" bson:"-"`
	BsonID  primitive.ObjectID `bson:"_id" json:"-"`
	User    *User              `json:"user" bson:"-"`
	MyMaids []string           `json:"my_maids"`
}

// Maid
type Maid struct {
	ID            string             `json:"id" bson:"-"`
	BsonID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Phone         string             `json:"phone"`
	Address       string             `json:"address"`
	Bio           string             `json:"bio"`
	ProfileImages []string           `json:"profile_images"`
	Carrers       []string           `json:"carrers"`
	User          *User              `json:"user" bson:"-" `
	CreatedBy     string             `json:"created_by"`
	Works         []*Work            `json:"works"`
	RatedBy       []string           `json:"rated_by" bson:"rated_by"` // This variable holds the data of those who rated this maid
	Rates         float32            `json:"rates"`
	RateCount     int                `json:"rate_count"`
}
type MaidUpdate struct {
	ID      string `json:"id"  bson:"-"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Bio     string `json:"bio"`
}

// Admin represents the admin  of the app
type Admin struct {
	ID     string             `json:"id"  bson:"-" `
	User   *User              `json:"user" bson:"-" `
	BsonID primitive.ObjectID `json:"-"  bson:"_id,omitempty"`
	// UserID    string `json:"user_id" bson:"user_id"`
	CreatedBy string `json:"created_by"`
}

// Work representnig the work and the shift the maid is available.
// the instance of this struct is to be stored in the database.
type Work struct {
	NO         uint
	Shift      int
	Type       int
	Experiance string
	Experties  []string
}

// RespWork ...
type RespWork struct {
	Shift      string   `json:"shift"`
	Type       string   `json:"type"`
	Experiance string   `json:"experiance"`
	Experties  []string `json:"experties"`
}

/*
	what is needed from maids
	- baby sitter
	- Permanent Maid Cooking
	- Part time Maid Cooker
	- Permanent Janitor
	-  Part time Janitor
	- Part time Loundary
	- Full time Loundary
	- Full time home care
	- Part time Home Care
	- All Inclusive

*/
