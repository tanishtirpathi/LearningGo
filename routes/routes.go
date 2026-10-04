package routes

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"go-task-api/controller"
)

// NewRouter registers the task endpoints and returns a handler for the HTTP server.
func NewRouter(collection *mongo.Collection) http.Handler {
	mux := http.NewServeMux()
	createTask := controller.CreateTask(collection)
	getTasks := controller.GetTasks(collection)
	deleteTask := controller.DeleteTask(collection)

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createTask(w, r)
		case http.MethodGet:
			getTasks(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

// Route keeps the original function name available to existing callers.
func Route(collection *mongo.Collection) http.Handler {
	return NewRouter(collection)
}
