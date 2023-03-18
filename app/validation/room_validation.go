package validation

type CreateRoomRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetRoomResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
