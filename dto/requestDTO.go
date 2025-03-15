package dto

type LoginRequest struct {
	Code string `json:"code"`
}

type GetUserInfoRequest struct {
	ID string `json:"id"`
}
