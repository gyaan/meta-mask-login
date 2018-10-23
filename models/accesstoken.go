package models

type AccessToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
