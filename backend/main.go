package main

import (
	"fmt"
	"log"
	"os"

	// "github.com/gin-gonic/gin"
	"wechat/database"
	"wechat/models"
	"wechat/routes"
	"wechat/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning .env file not loaded")
		return
	}
}

func main() {
	fmt.Println("helo world")

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	fmt.Println("Succefully connected")
	if err := db.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{}, &models.UserRoom{}); err != nil {
		log.Fatalf("Failed to migrate %v", err)
	}

	if err := db.SetupJoinTable(&models.Room{}, "Members", &models.UserRoom{}); err != nil {
		log.Fatalf("Failed to setup join table: %v", err)
	}

	websocket.InitWebSocketDB(db)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(router, db)
	router.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler(ctx.Writer, ctx.Request)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("✅ Server running on :%s\n", port)

	fmt.Println("🚀 Server started on http://localhost:8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
