package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eliday15/go_back/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Carga las variables de entorno desde el archivo .env
	godotenv.Load(".env")

	// Obtiene el valor del puerto desde las variables de entorno
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("$DB_URL must be set")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	fmt.Println("Port is: ", portString)

	// Crea un nuevo router Chi
	router := chi.NewRouter()

	// Configura el middleware CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Crea un subrouter para la versión 1 de la API
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/error", errorHandler)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFollow))
	v1Router.Get("/follows", apiCfg.middlewareAuth(apiCfg.handlerGetFollows))
	v1Router.Delete("/follows/{feed_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFollow))
	router.Mount("/v1", v1Router)

	// Configura el servidor HTTP
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Inicia el servidor
	log.Println("Server is running on port: ", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
