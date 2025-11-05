package models

import "time"

type Cat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CatTypeID int       `json:"cat_type_id"`
	MasterID  int       `json:"master_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	CatType *Type   `json:"cat_type,omitempty"`
	Master  *Master `json:"master,omitempty"`
}

type Type struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Master struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Place     string    `json:"place"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
