package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlistapp/internal/config"
	"wishlistapp/internal/database"
	"wishlistapp/internal/handler"
	"wishlistapp/internal/repository"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	wishRepo := repository.NewWishRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	wishlistHandler := handler.NewWishlistHandler(wishRepo)
	authHandler := handler.NewAuthHandler(userRepo)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/wishlist", wishlistHandler.GetWishlist)
	r.Post("/wishlist", wishlistHandler.CreateWish)
	r.Patch("/wishlist/{id}", wishlistHandler.UpdateWish)
	r.Delete("/wishlist/{id}", wishlistHandler.DeleteWish)
	r.Post("/auth/login", authHandler.Login)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
