package models

import "time"

type Capsula struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OpenDate  time.Time `db:"open_date" json:"open_date"`
	Sent      bool      `db:"sent" json:"sent"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateCapsulaReq struct {
	Name     string `json:"name"`
	OpenDate string `json:"open_date" example:"2025-12-10"`
}
