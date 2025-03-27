package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jt3angga/test-chatbot-be/internal/client"
	"github.com/jt3angga/test-chatbot-be/internal/handler"
	"github.com/jt3angga/test-chatbot-be/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    router := gin.Default()

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: false,
		}))
		
    router.Use(middleware.Logger())
    router.Use(middleware.Auth())
    router.Use(middleware.RateLimiter(5))

		apiKey := os.Getenv("GROQ_API_KEY")
    baseURL := os.Getenv("GROQ_BASE_URL")

    if apiKey == "" || baseURL == "" {
        log.Fatal("GROQ_API_KEY or GROQ_BASE_URL is not set")
    }

    groqClient := client.NewGroqClient(apiKey, baseURL)
    chatHandler := handler.NewChatHandler(groqClient)

    router.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
    pprof.Register(router)

    router.POST("/chat", chatHandler.ChatStream)
    router.Run(":8080")
}