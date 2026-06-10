package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"goalstack/internal/handlers"
	"goalstack/internal/storage"
)

func main() {
	// Initialize in-memory store
	store := storage.NewMemoryStore()

	// Initialize handler
	handler := handlers.NewHandler(store)

	// Create Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(corsMiddleware())

	// Goal routes
	router.POST("/goals", handler.CreateGoal)
	router.GET("/goals", handler.ListGoals)
	router.GET("/goals/:id", handler.GetGoal)
	router.PUT("/goals/:id", handler.UpdateGoal)
	router.DELETE("/goals/:id", handler.DeleteGoal)

	// Subtask routes
	router.POST("/goals/:id/subtasks", handler.AddSubtask)
	router.PUT("/goals/:id/subtasks/:taskId", handler.UpdateSubtask)
	router.PATCH("/goals/:id/subtasks/reorder", handler.ReorderSubtasks)
	router.DELETE("/goals/:id/subtasks/:taskId", handler.DeleteSubtask)
	router.PATCH("/subtasks/:taskId/progress", handler.UpdateSubtaskProgress)
	router.PATCH("/subtasks/:taskId/complete", handler.CompleteSubtask)

	// Note routes
	router.POST("/goals/:id/subtasks/:taskId/notes", handler.AddNote)
	router.DELETE("/goals/:id/subtasks/:taskId/notes/:noteId", handler.DeleteNote)

	// Link routes
	router.POST("/goals/:id/subtasks/:taskId/links", handler.AddLink)
	router.DELETE("/goals/:id/subtasks/:taskId/links/:linkId", handler.DeleteLink)

	// Checklist routes
	router.POST("/goals/:id/subtasks/:taskId/checklist", handler.AddChecklistItem)
	router.PATCH("/goals/:id/subtasks/:taskId/checklist/:itemId", handler.UpdateChecklistItem)
	router.DELETE("/goals/:id/subtasks/:taskId/checklist/:itemId", handler.DeleteChecklistItem)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// corsMiddleware returns a CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
