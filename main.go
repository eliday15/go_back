package main

import(
	"fmt"
	"os"
	"log"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
)
func main()  {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}
	fmt.Println("Port is: ", portString)
	
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/error", errorHandler)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}
	

	log.Println("Server is running on port: ", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}