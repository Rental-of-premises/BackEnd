package api

import (
	"fmt"
	"log"
	"net/http"
	api_controllers "rent/internal/api/controllers"
	"rent/internal/config"
	"rent/internal/storage/repository"

	"github.com/gorilla/mux"
)

func CreateAndRunRoutes() {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})

	userController := &api_controllers.UserController{repository.GetUserRepository()}

	r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET") 
	r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST")
	r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST")

	port := config.GetSingletonConfig().ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
