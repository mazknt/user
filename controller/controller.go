package controller

import (
	"authentication/dto"
	"authentication/service"
	json_util "authentication/util/json"
	"encoding/json"
	"log"
	"net/http"

	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
)

type Controller struct {
	Service service.ServiceInterface
}

func NewController(Service service.Service) *Controller {
	return &Controller{
		Service: Service,
	}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	userInfo := FP.Pipe3(
		json_util.ReadRequest[dto.LoginRequest](r),
		E.Map[error](func(req dto.LoginRequest) string { return req.Code }),
		c.Service.Login,
		E.Fold(
			func(err error) dto.LoginResponse {
				log.Println("error: %w", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				return dto.LoginResponse{}
			},
			func(res dto.LoginResponse) dto.LoginResponse {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				return res
			},
		),
	)
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Println("Error encoding userInfo:", err)
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}
	return
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	userInfo := FP.Pipe3(
		json_util.ReadRequest[dto.GetUserInfoRequest](r),
		E.Map[error](func(req dto.GetUserInfoRequest) string { return req.ID }),
		c.Service.GetUser,
		E.Fold(
			func(err error) dto.GetUserInfoResponse {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				return dto.GetUserInfoResponse{}
			},
			func(res dto.GetUserInfoResponse) dto.GetUserInfoResponse {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				return res
			},
		),
	)

	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Println("Error encoding userInfo:", err)
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}

	return
}
