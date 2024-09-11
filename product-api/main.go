package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	protos "github.com/martbul/currency/protos/currency"
	"github.com/martbul/product-api/data"
	"github.com/martbul/product-api/handlers"
	"google.golang.org/grpc"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	
	l := hclog.Default( )
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	//create gRPC client
	cc := protos.NewCurrencyClient(conn)

	//create products db
	db := data.NewProductsDB(cc, l)

	productsHandler := handlers.NewProducts(l, v, db)

	

	// adding gorila serveMux(allowes for registering more detailed routers)
	serveMux := mux.NewRouter()


	// handlers for API
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productsHandler.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products", productsHandler.ListAll)
	
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.ListSingle).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.ListSingle)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productsHandler.Create)
	postRouter.Use(productsHandler.MiddlewareValidateProduct)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9+]}", productsHandler.Update)
	putRouter.Use(productsHandler.MiddlewareValidateProduct)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", productsHandler.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	swaggerHandler := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", swaggerHandler)
	//serving a the swagger.yaml file
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// create new server
	// server := &http.Server{
	// 	Addr: ":9090",
	// 	// using the new created serverMux, instead of the default
	// 	Handler:      ch(serveMux),
	// 	IdleTimeout:  120 * time.Second,
	// 	ReadTimeout:  1 * time.Second,
	// 	WriteTimeout: 1 * time.Second,
	// }

	server := http.Server{
		Addr:         ":9090",                                     // configure the bind address
		Handler:      ch(serveMux),                                           // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	//! DON`T UNDERSTAND
	//wrappingt he service in a go func in order to not block
	go func() {
		l.Info("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	// broadcasting a message on the sigChan whenever an opperating system kill's command or interupt(now when you do ctrl + c and kill the running server it will gracefuly shutdown)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Got signal:", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second) //allowing 30 sec for gracefuls shutdow, after them the server will forcefully shutdown

	// this is graceful shutdown,the server will no longer accept new requests and will wait until it has completed all the old requests, before shuting down
	server.Shutdown(timeoutContext)
}
