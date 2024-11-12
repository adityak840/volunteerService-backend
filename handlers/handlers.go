package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/volunteerService-backend/services"
)

var todo services.Todo

// Response struct to standardize all responses
type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// healthCheck - simple function to test api if its working
func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Health Check",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := todo.GetAllTodos()
	if err != nil {
		log.Println(err)
		res := Response{
			Msg:  "Error retrieving todos",
			Code: 500,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todos)
}

func getTodoById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	todo, err := todo.GetTodoById(id)
	if err != nil {
		log.Println(err)
		res := Response{
			Msg:  "Todo not found",
			Code: 404,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todo)
}

func getTodoByVol(w http.ResponseWriter, r *http.Request) {
	// Retrieve the 'volType' query parameter from the URL
	volType := r.URL.Query().Get("volType")
	if volType == "" {
		// If no 'volType' is provided, return an error response
		errorRes := Response{
			Msg:  "Missing 'volType' query parameter",
			Code: 400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	// Call the service method to get todos by VolunteerType
	todos, err := todo.GetTodosByVolType(volType)
	if err != nil {
		log.Println("Error retrieving todos by volunteer type:", err)
		errorRes := Response{
			Msg:  "Error retrieving todos",
			Code: 500,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	// Return the list of todos as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todos)
}

func getTodoByOrg(w http.ResponseWriter, r *http.Request) {
	// Retrieve the 'orgName' query parameter from the URL
	orgName := r.URL.Query().Get("orgName")
	if orgName == "" {
		// If no 'orgName' is provided, return an error response
		errorRes := Response{
			Msg:  "Missing 'orgName' query parameter",
			Code: 400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	// Call the service method to get todos by OrganisationName
	todos, err := todo.GetTodosByOrg(orgName)
	if err != nil {
		log.Println("Error retrieving todos by organisation name:", err)
		errorRes := Response{
			Msg:  "Error retrieving todos",
			Code: 500,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	// Return the list of todos as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatal(err)
	}

	err = todo.InsertTodo(todo)
	if err != nil {
		errorRes := Response{
			Msg:  "Error creating todo",
			Code: 304,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	res := Response{
		Msg:  "Successfully created todo",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		res := Response{
			Msg:  "Error decoding request",
			Code: 400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = todo.UpdateTodo(id, todo)
	if err != nil {
		errorRes := Response{
			Msg:  err.Error(),
			Code: 500,
		}
		jsonStr, err := json.Marshal(errorRes)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		w.Write(jsonStr)
		return
	}

	res := Response{
		Msg:  "Successfully updated todo",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := todo.DeleteTodo(id)
	if err != nil {
		errorRes := Response{
			Msg:  "Error deleting todo",
			Code: 304,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	res := Response{
		Msg:  "Successfully deleted todo",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

// SignupHandler handles the signup request
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user services.User

	// Decode the request body into user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the Signup function
	userType, err := services.Signup(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Send back the userType as the response
	response := struct {
		UserType string `json:"userType"`
	}{
		UserType: userType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		sendErrorResponse(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

// LoginHandler handles the login request
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body into loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the Login function
	_, err = services.Login(w, loginRequest.Email, loginRequest.Password)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		sendErrorResponse(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	users, err := services.GetUsersByID(ids)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
