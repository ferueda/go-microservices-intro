package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ferueda/go-microservices-intro/models"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := models.GetProducts()

	w.Header().Set("content-type", "application/json")
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(models.Product)

	models.AddProduct(prod)
	w.WriteHeader(http.StatusCreated)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var prod models.Product
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http.Error(w, "Id badformatted", http.StatusBadRequest)
		return
	}

	err = models.UpdateProduct(id, prod)
	if err != nil {
		if err == models.ErrProductNotFound {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unknown error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prod models.Product
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		if err := prod.Validate(); err != nil {
			http.Error(w, "Bad product format", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
