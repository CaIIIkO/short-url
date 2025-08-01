package main

import (
	"log"
	"net/http"
	"os"
	"short-url/internal/auth"
	"short-url/internal/db"
	"short-url/internal/user"
)

func main() {
	log.Println("short-url-api is starting...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "supersecretjwtkey"
	}

	dsn := os.Getenv("DATABASE_DSN")

	//Подключение к базе данных
	pool := db.Connect(dsn)

	//Инциализация jwtManager
	jwtManager := auth.NewJWTManager(secretKey)

	//инциализация хэндлеров и сервисов
	userRepo := user.NewUserRepository(pool)
	userService := user.NewUserService(userRepo, jwtManager)
	userHandler := user.NewUserHandler(userService)

	//http
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", userHandler.Register) //POST
	mux.HandleFunc("/auth/login", userHandler.Login)       //POST

	// mux.HandleFunc("/url/{code}")       //GET
	// mux.HandleFunc("/url/{code}/stats") //GET
	// mux.HandleFunc("/url/create")       //POST

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Запуск сервера
	log.Printf("server is running http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
