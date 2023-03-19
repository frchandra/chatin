package validation

import "time"

type CreateRoomRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetRoomResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
