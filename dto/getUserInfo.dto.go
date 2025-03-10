package dto

type GetUserInfoRequest struct {
	ID string `json:"id"`
}

type GetUserInfoResponse struct {
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Picture string   `json:"picture"`
	Friends []string `json:"friends"`
}
