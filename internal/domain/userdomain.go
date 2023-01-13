package domain

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
