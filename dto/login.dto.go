package dto

type LoginRequest struct {
	Code string `json:"code"`
}

type LoginResponse struct {
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Picture string   `json:"picture"`
	Friends []string `json:"friends"`
}
