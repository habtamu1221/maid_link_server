package model

type Success struct {
	Success bool `json:"success"`
}

type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type LoginSuccess struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
