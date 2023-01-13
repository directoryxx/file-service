package domain

type File struct {
	ID     int    `json:"-"`
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	UserID int    `json:"user_id"`
	IsTemp int    `json:"is_temp"`
}
