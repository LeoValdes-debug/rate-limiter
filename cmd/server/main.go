package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/leovaldes-debug/rate-limiter/internal/handler"
	"github.com/leovaldes-debug/rate-limiter/internal/limiter"
	"github.com/leovaldes-debug/rate-limiter/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	capacity := getEnvFloat("RATE_CAPACITY", 10)
	refill := getEnvFloat("RATE_REFILL", 5)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	l := limiter.New(capacity, refill)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)
	mux.HandleFunc("/hello", handler.Hello)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      middleware.RateLimit(l)(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s (capacity=%.0f refill=%.0f/s)", port, capacity, refill)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("stopped")
}

func getEnvFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}
