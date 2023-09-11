package models

type Login struct {
	Id int64 `json:"id,omitempty"`

	Password string `json:"password,omitempty"`

	Username string `json:"username,omitempty"`
}
