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

type AdminInfo struct {
	Success bool   `json:"success"`
	Admin   *Admin `json:"admin"`
}

type UserInfo struct {
	Success bool        `json:"success"`
	User    interface{} `json:"user"`
}
type Profile struct {
	Success bool        `json:"success"`
	Role    int         `json:"role"`
	Body    interface{} `json:"body"`
}

// PayForMaid this struct will beused when the client pays for the information
// about the maid ..... the acount number and the  password are placed in the accounts.csv file
// with their amount of money they have
type PayForMaid struct {
	AccNO    string `json:"account"`
	Password string `json:"password"`
	MaidID   string `json:"maid_id"`
}
