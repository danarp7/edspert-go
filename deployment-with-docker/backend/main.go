package main

import (
	"log"
	"os"
	"postgres/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize gin
	r := gin.Default()

	// Get the config from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load db config
	db, err := config.OpenDB(os.Getenv("POSTGRES_URL"), true)
	if err != nil {
		log.Fatalln(err)
	}
	defer config.CloseDB(db)

	// Init clean arch
	repository := config.InitRepository(db)
	usecase := config.InitUsecase(repository.AlbumRepository)
	handler := config.InitHandler(usecase.AlbumUsecase)

	// Create the API
	albumRoutes := r.Group("/api/v1/albums")
	{
		albumRoutes.GET("/", handler.AlbumHandler.GetAllAlbum)
		albumRoutes.POST("/", handler.AlbumHandler.Create)
		albumRoutes.POST("/batch", handler.AlbumHandler.BatchCreate)
		albumRoutes.GET("/:id", handler.AlbumHandler.Get)
		albumRoutes.PUT("/:id", handler.AlbumHandler.Update)
		albumRoutes.DELETE("/:id", handler.AlbumHandler.Delete)
	}

	// Run the gin gonic in port 5000
	r.Run(":4000")
}
