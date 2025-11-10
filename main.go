package main

import (
	"todo/database"
	"todo/handlers"
	"todo/middleware"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()
	// routes
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", handlers.RegisterUser)
	r.HandleFunc("/login", handlers.Login)
	

	// protected routes
	protectRoutes := r.PathPrefix("/").Subrouter()
	protectRoutes.Use(middleware.AuthMiddleware)

	protectRoutes.HandleFunc("/create", handlers.CreateJob)
	
	protectRoutes.HandleFunc("/update", handlers.UpdateTodo)
	protectRoutes.HandleFunc("/delete", handlers.DeleteTodo) 
    protectRoutes.HandleFunc("/todo", handlers.GetAllTodo)
	// start the server
	
	//fmt.Println("Server is running")
	//err := http.ListenAndServe(":8081", r)
	//if err != nil {
		//panic(err)
		port := os.Getenv("PORT")
     if port == "" {
      port = "8081" 
    }

     fmt.Println("Server is running on port " + port)
      err := http.ListenAndServe(":"+port, r)
     if err != nil {
     panic(err)
}
}
	


