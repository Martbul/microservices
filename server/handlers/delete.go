package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/martbul/microservices/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of pr
// responses:
// 201: noContent

//DeleteProduct deletes a produc from the database
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])

	p.l.Println("Handler DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product jnot found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
	
}