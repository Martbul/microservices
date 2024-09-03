// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/martbul/microservices/data"
)

// A list of products returns in the response
// swagger:response productResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameter deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
// 	listOfProducts := data.GetProducts()
// 	err := listOfProducts.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
// 	}
// }

// func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
// 	p.l.Println("Handle POST product")

// 	prod := r.Context().Value(KeyProduct{}).(data.Product)
// 	data.AddProduct(&prod)
// }

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	//gorilla is extracting the varibles from the url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id to int", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT product", id)
	//! check
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// validating thr request
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//validate the product
		err = prod.Validate()
		if err != nil {

			p.l.Println("[ERROR] validationg product", err)
			http.Error(rw, fmt.Sprintf("Error validation product: %s",err), http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)

	})
}
