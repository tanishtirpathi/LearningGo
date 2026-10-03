package routes

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	Completed bool   `json:"completed" bson:"completed"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	result, err := tasksCollection.InsertOne(r.Context(), task)

	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	task.ID = result.InsertedID.(string)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
func GetTasks(w http.ResponseWriter, r *http.Request) {
	cursor, err := tasksCollection.Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(r.Context())
	var tasks []Task

	err = cursor.All(r.Context(), &tasks)
	if err != nil {
		http.Error(w, "Failed to decode tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
