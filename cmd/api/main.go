package main

import (
	"log"
	"os"

	"github.com/PatrikMaltacm/life-uptime/internal/database"
	"github.com/PatrikMaltacm/life-uptime/internal/handler"
	"github.com/PatrikMaltacm/life-uptime/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	dsn := os.Getenv("DATABASE_URL")
	db := database.Connect(dsn)
	defer db.Close()

	monitorRepo := repository.NewMonitorRepository(db)

	monitorHandler := handler.NewMonitorHandler(monitorRepo)

	router := gin.Default()
	v1 := router.Group("/api/v1")

	handler.InitRoutes(v1, monitorHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
