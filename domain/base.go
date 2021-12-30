package domain

import "time"

type Base struct {
	ID        string     `json:"id" bson:"_id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:",omitempty"`
}
