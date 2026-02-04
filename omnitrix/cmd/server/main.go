package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/llm"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/provider"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/reflex"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/router"
)

// Comments below for Swagger

// @title           Omnitrix AI Gateway
// @version         1.0
// @description     A High-Performance Intent-Based Router for LLMs.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@omnitrix.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	port := ":8080"

	// Initialize reflex engine
	reflexEngine := reflex.New()

	// Ollama Client
	ollamaClient := llm.NewOllamaClient("http://localhost:11434", "phi3:mini") // "phi3:mini" is default for classifier

	// Groq Client
	groqKey := os.Getenv("GROQ_API_KEY")
	var groqClient *llm.GroqClient
	if groqKey != "" {
		groqClient = llm.NewGroqClient(groqKey)
		log.Println("Groq Client Initialized")
	} else {
		log.Println("Warning: GROQ_API_KEY not found. Premium models will fail.")
	}

	// Client Factory
	providerFactory := provider.NewFactory(ollamaClient, groqClient)

	// Router
	r := router.New(reflexEngine, ollamaClient, providerFactory)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		log.Printf("Omnitrix Gateway starting on %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful Shutdown Logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
