package models

type Profile struct {
	Id int64 `json:"id,omitempty"`

	Address string `json:"address,omitempty"`

	Age int64 `json:"age,omitempty"`

	Name string `json:"name,omitempty"`
}
