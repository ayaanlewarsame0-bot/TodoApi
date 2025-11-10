package handlers

import (
	"todo/database"
	"todo/middleware"
	"todo/models"
	"todo/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	

)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// collect the details of the user as request body
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the user already exist
	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	// Hash the password
	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusBadRequest)
		return
	}

	req.Password = HashPassword

	// add the user to the database
	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "unable to create user", http.StatusBadRequest)
		return
	}

	// send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User created successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	// decode the request the request body
	var login models.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the user exists
	var user models.User
	err = database.Db.Where("email = ?", login.Email).First(&user).Error
	if err != nil {
		http.Error(w, "this user does not exist", http.StatusBadRequest)
		return
	}

	// check if password matches what we have in our database
	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	// uint ---> int ---> string.

	idStr := strconv.Itoa(int(user.ID))

	// generating a token
	token, err := middleware.GenerateJWT(idStr)
	if err != nil {
		http.Error(w, "unable to generate token", http.StatusInternalServerError)
		return
	}

	// send a response
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var post models.Todo
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.UserID = userID

	err = database.Db.Create(&post).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Post created successfully")
}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Todo
	err := database.Db.Find(&jobs).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(jobs)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract todo ID from query params (e.g. /update?id=1)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing todo id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	// Decode updated todo info
	var updatedJob models.Todo
	err = json.NewDecoder(r.Body).Decode(&updatedJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify user ID from token
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Find the existing todo
	var job models.Todo
	err = database.Db.First(&job, id).Error
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	// Ensure the todo belongs to the logged-in user
	if job.UserID != userID {
		http.Error(w, "unauthorized to update this todo", http.StatusForbidden)
		return
	}

	
	if updatedJob.Title != "" {
		job.Title = updatedJob.Title
	}
	if updatedJob.Description != "" {
		job.Description = updatedJob.Description
	}
	if updatedJob.Status != "" {
		job.Status = updatedJob.Status
	}

	err = database.Db.Save(&job).Error
	if err != nil {
		http.Error(w, "unable to update todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "todo updated successfully")
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract todo ID from query params (e.g. /delete?id=1)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing todo id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	// Verify user ID from token
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Find the Todo
	var job models.Todo
	err = database.Db.First(&job, id).Error
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	// Check if it exists
	if job.UserID != userID {
		http.Error(w, "unauthorized to delete this todo", http.StatusForbidden)
		return
	}

	// Delete Todo
	err = database.Db.Delete(&job).Error
	if err != nil {
		http.Error(w, "unable to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "todo deleted successfully")
}

	


