package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	config "github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	viper "github.com/spf13/viper"
)

func main() {
	var err error
	viper.SetDefault("SERVER_TIMEZONE", "UTC")
	loc, err := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	if err != nil {
		log.Fatalf("Invalid timezone: %v", err)
	}
	time.Local = loc

	/** Init config */
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig fail: %v", err)
	}

	log.Println(env.App)
	log.Println(env.Database)
	log.Println(env.Jwt)

	// r := gin.Default()

	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	} else {
		log.Println("Server exited properly")
	}
}
