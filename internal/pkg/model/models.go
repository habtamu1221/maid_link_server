package model

// User represents all the users in the system
type User struct {
	ID       string `json:"id" bson:"_id,omitempty" `
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImageUrl string `json:"image_url"`
	Role     int8   `json:"role"`
}

// Client represents the client
type Client struct {
	ID      string   `json:"id"`
	User    *User    `json:"user"`
	MyMaids []string `json:"my_maids"`
}

// Maid
type Maid struct {
	ID            string   `json:"id" bson:"_id,omitempty"`
	Phone         string   `json:"phone"`
	Address       string   `json:"address"`
	Bio           string   `json:"bio"`
	ProfileImages []string `json:"profile_images"`
	Carrers       []string `json:"carrers"`
	User          *User    `json:"user"`
	CreatedBy     string   `json:"created_by"`
	Works         []*Work  `json:"works"`
	Rates         int      `json:"rates"`
	RateCount     int      `json:"rate_count"`
}

// Admin represents the admin  of the app
type Admin struct {
	ID        string `json:"id"  bson:"_id,omitempty" `
	User      *User  `json:"user" bson:"-" `
	UserID    string `json:"user_id" bson:"user_id"`
	CreatedBy string `json:"created_by"`
}

// Work representnig the work and the shift the maid is available.
type Work struct {
	Shift int
	Type  int
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
