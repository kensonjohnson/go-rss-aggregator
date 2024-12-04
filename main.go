package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/kensonjohnson/go-rss-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading environment variables")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN environment variable not set")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handlerReadiness)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	v1Router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Delete("/feed-follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Println("Server listening on port " + portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
