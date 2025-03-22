package controller

import (
	email "authentication/Domain/models/User/Email"
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
	Service service.UserServiceInterface
}

func NewController(Service service.UserServiceInterface) *Controller {
	return &Controller{
		Service: Service,
	}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	userInfo := FP.Pipe3(
		json_util.ReadRequest[dto.LoginRequest](r),
		E.Map[error](func(req dto.LoginRequest) string { return req.Code }),
		E.Chain(c.Service.Login),
		E.Fold(
			func(err error) dto.UserInformation {
				log.Println("error: %w", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				return dto.UserInformation{}
			},
			func(res dto.UserInformation) dto.UserInformation {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				return res
			},
		),
	)

	if err := json.NewEncoder(w).Encode(dto.CreateUserResponse(userInfo)); err != nil {
		log.Println("Error encoding userInfo:", err)
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}
	return
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	userInfo := FP.Pipe4(
		json_util.ReadRequest[dto.GetUserInfoRequest](r),
		E.Map[error](func(req dto.GetUserInfoRequest) string { return req.ID }),
		E.Chain(func(id string) E.Either[error, email.Email] { return email.NewEmail(id) }),
		E.Chain(c.Service.GetUser),
		E.Fold(
			func(err error) dto.UserInformation {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				return dto.UserInformation{}
			},
			func(res dto.UserInformation) dto.UserInformation {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				return res
			},
		),
	)

	if err := json.NewEncoder(w).Encode(dto.CreateUserResponse(userInfo)); err != nil {
		log.Println("Error encoding userInfo:", err)
		http.Error(w, "Failed to encode user info", http.StatusInternalServerError)
	}

	return
}
