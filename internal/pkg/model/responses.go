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

type PasswordChangeSuccess struct {
	NewPassword string `json:"new_password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ShortError struct {
	Err string `json:"err"`
}

type ShortSuccess struct {
	Msg string `json:"msg"`
}
