package models

import "time"

type Message struct {
	ID        string    `db:"id" json:"id"`
	Message   string    `db:"message" json:"message"`
	PhotoURL  string    `db:"photo_url" json:"photo_url"`
	CapsulaID string    `db:"capsula_id" json:"capsula_id"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateMessageReq struct {
	CapsulaID string `json:"-"`
	Message   string `json:"message"`
	Email     string `json:"email"`
	PhotoURL  string `json:"photo_url"`
}
