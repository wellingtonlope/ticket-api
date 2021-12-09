package domain

import "time"

type Base struct {
	ID        string     `json:"id" bson:"_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:",omitempty"`
}
