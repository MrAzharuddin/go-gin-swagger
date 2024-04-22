package models

type Employee struct {
	Id int64 `json:"id,omitempty"`

	Age int `json:"age,omitempty"`

	Company string `json:"company,omitempty"`

	Name string `json:"name,omitempty"`
}
