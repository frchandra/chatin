package validation

type GetClientResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
