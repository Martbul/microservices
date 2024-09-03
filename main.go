package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/martbul/microservices/data"
	"github.com/martbul/microservices/handlers"
)

func main() {
	
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	productsHandler := handlers.NewProducts(l, v)

	// adding gorila serveMux(allowes for registering more detailed routers)
	serveMux := mux.NewRouter()

	// handlers for API
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productsHandler.GetProducts)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productsHandler.Create)
	postRouter.Use(productsHandler.MiddlewareValidateProduct)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9+]}", productsHandler.Update)
	putRouter.Use(productsHandler.MiddlewareValidateProduct)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", productsHandler.DeleteProduct)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	swaggerHandler := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", swaggerHandler)
	//serving a the swagger.yaml file
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create new server
	server := &http.Server{
		Addr: ":9090",
		// using the new created serverMux, instead of the default
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//! DON`T UNDERSTAND
	//wrappingt he service in a go func in order to not block
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	// broadcasting a message on the sigChan whenever an opperating system kill's command or interupt(now when you do ctrl + c and kill the running server it will gracefuly shutdown)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second) //allowing 30 sec for gracefuls shutdow, after them the server will forcefully shutdown

	// this is graceful shutdown,the server will no longer accept new requests and will wait until it has completed all the old requests, before shuting down
	server.Shutdown(timeoutContext)
}
