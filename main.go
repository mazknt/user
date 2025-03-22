package main

import (
	"authentication/controller"
	"authentication/interface/api"
	"authentication/interface/repository"
	"authentication/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	ctrl := controller.NewController(
		service.NewService(
			// api.GoogleAPI{},
			api.GoogleAPI{},
			repository.FireStore{},
		),
	)

	r := mux.NewRouter()
	r.HandleFunc("/login", ctrl.Login).Methods("POST")
	r.HandleFunc("/get-user-info", ctrl.GetUser).Methods("POST")

	// CORS設定（全てのドメインからのリクエストを許可）
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // 任意のオリジンを許可（セキュリティリスクを避けるため本番環境では必要に応じてドメインを指定）
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
