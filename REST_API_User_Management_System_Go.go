package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.RWMutex
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if user.Name == "" || user.Email == "" {
			http.Error(w, "Name and email required", http.StatusBadRequest)
			return
		}

		mu.Lock()
		user.ID = nextID
		nextID++
		users[user.ID] = user
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()

		json.NewEncoder(w).Encode(users)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		mu.RLock()
		user, exists := users[id]
		mu.RUnlock()

		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var updated User

		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if updated.Name == "" || updated.Email == "" {
			http.Error(w, "Name and email required", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if _, exists := users[id]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		updated.ID = id
		users[id] = updated
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()

		if _, exists := users[id]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		delete(users, id)
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/users/", userByIDHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	fmt.Println("Server running on port 8080")

	log.Fatal(server.ListenAndServe())
}
