package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ferueda/go-microservices-intro/models"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		{
			reg := regexp.MustCompile(`/([0-9]+)`)
			g := reg.FindAllStringSubmatch(r.URL.Path, -1)

			if len(g) != 1 || len(g[0]) != 2 {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(g[0][1])
			if err != nil {
				http.Error(w, "Invalid product id", http.StatusBadRequest)
				return
			}

			p.updateProduct(id, w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := models.GetProducts()

	w.Header().Set("content-type", "application/json")
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	var prod models.Product
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	models.AddProduct(prod)
	w.WriteHeader(http.StatusCreated)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	var prod models.Product
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err := models.UpdateProduct(id, prod)
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
