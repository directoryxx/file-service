package domain

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type UserResponseAuthService struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Profile User   `json:"profile"`
}

type PublishAuthLogout struct {
	Data   LogoutAction
	Action string
}

type PublishAuthLogin struct {
	Data   LoginAction
	Action string
}

type LogoutAction struct {
	Uuid string
}

type LoginAction struct {
	Uuid string
	User User
	Exp  time.Duration
}
