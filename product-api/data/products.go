//! making a call to the  EuropeCentralBank for the currency rates. The rates dont change very often so making new request for every item is not efficient. 
//! Cause of that i will use gRPC for bi-directional straming and pull architecture(the ECB api will notify product-api if there is a change and then a new call will be made)
package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	protos "github.com/martbul/currency/protos/currency"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{currency: c, log: l}
}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	p.log.Info("currency",currency)

	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	p.log.Info("rate",rate)

	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	products := Products{}
	for _, p := range productList {
		newProduct := *p
		newProduct.Price = newProduct.Price * rate
		products = append(products, &newProduct)
	}

	return products, nil
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)

	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}


	 newProduct := *productList[i]
	 newProduct.Price = newProduct.Price * rate

	return &newProduct, nil
}
// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(pr.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &pr

	return nil
}

// AddProduct adds a new product to the database
func (p *ProductsDB) AddProduct(pr Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	pr.ID = maxID + 1
	productList = append(productList, &pr)
}

// DeleteProduct deletes a product from the database
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	rateRequest  := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rateRequest )
	return resp.Rate, err
}


var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
