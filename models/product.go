package models

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"required,gte=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type Products []Product

func (pp *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(&pp)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(&p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}

func GetProducts() Products {
	return productList
}

func AddProduct(p Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func UpdateProduct(id int, p Product) error {
	p.ID = id
	idx, err := findProductIdx(id)
	if err != nil {
		return err
	}

	productList[idx] = p
	return nil
}

func getNextId() int {
	return productList[len(productList)-1].ID + 1
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProductIdx(id int) (int, error) {
	for i, p := range productList {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrProductNotFound
}

var productList = []Product{
	Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffe",
		Price:       2.45,
		SKU:         "abc123",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "Another description",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}
