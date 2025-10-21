package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

type Kennel struct {
	dogs map[string]Dog
	mu   sync.Mutex
}

func NewKennel() *Kennel {
	return &Kennel{
		dogs: make(map[string]Dog),
	}
}

func (k *Kennel) GetDog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dogID := r.PathValue("id")

		if dogID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		k.mu.Lock()
		defer k.mu.Unlock()
		dog, ok := k.dogs[dogID]

		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": "dog not found",
				"id":    dogID,
			})
			return
		}

		if err := json.NewEncoder(w).Encode(dog); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (k *Kennel) PostDog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dogID := uuid.New().String()

		var dog Dog
		if err := json.NewDecoder(r.Body).Decode(&dog); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		k.mu.Lock()
		defer k.mu.Unlock()
		k.dogs[dogID] = dog

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(dogID)
	}
}

func (k *Kennel) DeleteDog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dogID := r.PathValue("id")

		if dogID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		k.mu.Lock()
		defer k.mu.Unlock()

		delete(k.dogs, dogID)
		w.WriteHeader(http.StatusOK)
	}
}
