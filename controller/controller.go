package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Task struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Completed bool          `json:"completed" bson:"completed"`
}

// CreateTask validates and stores one task, then returns the stored document.
func CreateTask(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		err := json.NewDecoder(r.Body).Decode(&task)

		fmt.Println("Task received:", task)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		result, err := collection.InsertOne(r.Context(), task)

		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}
		task.ID, _ = result.InsertedID.(bson.ObjectID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(task)
	}
}

// GetTasks reads all tasks and returns them as a JSON array.
func GetTasks(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cursor, err := collection.Find(r.Context(), bson.D{})
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
		_ = json.NewEncoder(w).Encode(tasks)
	}
}

func DeleteTask(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		err := json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result, err := collection.DeleteOne(r.Context(), bson.M{"_id": task.ID})
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}
		if result.DeletedCount == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
