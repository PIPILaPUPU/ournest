package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlistapp/internal/dateideas"
	"wishlistapp/internal/notifications"
	"wishlistapp/internal/platform/auth"
	"wishlistapp/internal/platform/config"
	"wishlistapp/internal/platform/database"
	"wishlistapp/internal/wishlist"

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

	notificationRepo := notifications.NewRepository(pool, cfg.NotificationsEnabled)
	wishRepo := wishlist.NewWishRepository(pool, notificationRepo)
	userRepo := auth.NewUserRepository(pool)
	dateIdeasRepo := dateideas.NewRepository(pool, notificationRepo)
	tokenManager := auth.NewTokenManager(cfg.JWTSecret)
	wishlistHandler := wishlist.NewWishlistHandler(wishRepo)
	authHandler := auth.NewAuthHandler(userRepo, tokenManager)
	dateIdeasHandler := dateideas.NewHandler(dateIdeasRepo)
	notificationHandler := notifications.NewHandler(
		notificationRepo,
		cfg.WebPushPublicKey,
		cfg.NotificationsEnabled,
	)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Route("/auth", authHandler.Routes)
	r.Route("/wishlist", func(r chi.Router) {
		r.Use(tokenManager.Middleware)
		wishlistHandler.Routes(r)
	})
	r.Route("/date-ideas", func(r chi.Router) {
		r.Use(tokenManager.Middleware)
		dateIdeasHandler.Routes(r)
	})
	r.Get("/notifications/vapid-public-key", notificationHandler.PublicKey)
	r.Route("/notifications", func(r chi.Router) {
		r.Use(tokenManager.Middleware)
		notificationHandler.Routes(r)
	})

	workerCtx, stopWorker := context.WithCancel(context.Background())
	defer stopWorker()
	if cfg.NotificationsEnabled {
		sender := notifications.NewWebPushSender(
			cfg.WebPushPublicKey,
			cfg.WebPushPrivateKey,
			cfg.WebPushSubject,
		)
		worker := notifications.NewWorker(notificationRepo, sender, log.Default())
		go worker.Run(workerCtx)
	}

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
	stopWorker()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
