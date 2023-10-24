// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// -application/json
//
// Produces:
// -application/json
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aksentijevicd1/go-micros/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

/*func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT!")
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Printf("%q\n", reg.FindAllStringSubmatch(r.URL.Path, -1))
		p.l.Printf("%d\n", len(g[0]))
		p.l.Printf("%d\n", len(g[0][1]))
		if len(g) != 1 {
			p.l.Println("INVALID URI MORE THAN ONE ID")
			http.Error(w, "invalid URL", http.StatusNotFound)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("INVALID URI MORE THAN ONE CAPTURE GROUP")

			http.Error(w, "invalid URL", http.StatusNotFound)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Unable to convert to number", idString)
			http.Error(w, "INVALID URL", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
*/

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("PROD: %#v", prod)
	data.AddProduct(&prod)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Unable to parse to int", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct{}

//checks product in req and calls next if ok
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("error, deserializing product", err)
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//validate the product

		err = prod.Validate()
		if err != nil {
			p.l.Println("error, validating product", err)
			http.Error(
				w,
				fmt.Sprintf("Unable validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		//add product to context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		r = r.WithContext(ctx)
		//calls next handler which can be new middleware in chain
		next.ServeHTTP(w, r)
	})
}
