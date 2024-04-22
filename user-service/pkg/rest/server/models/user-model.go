package models

type User struct {
	Id int64 `json:"id,omitempty"`

	Address string `json:"address,omitempty"`

	Age int `json:"age,omitempty"`

	Email string `json:"email,omitempty"`

	Name string `json:"name,omitempty"`

	Phone string `json:"phone,omitempty"`
}
