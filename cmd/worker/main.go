package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/risefit/email-worker/pkg/middleware"
)

func main() {
	srv, err := InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	r := gin.Default()

	r.POST("/send-email", middleware.InternalAuth(srv.Config.InternalAPIKey), srv.Handler.SendEmail)

	log.Printf("Starting email worker service on port %s", srv.Config.Port)
	if err := r.Run(":" + srv.Config.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
