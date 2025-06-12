package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	config "github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/router"
	viper "github.com/spf13/viper"
)

func main() {
	var err error

	// Set up the timezone so when program executes, it knows what timezone it is in
	viper.SetDefault("SERVER_TIMEZONE", "UTC")
	loc, err := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	if err != nil {
		log.Fatalf("Invalid timezone: %v", err)
	}
	time.Local = loc

	// load the config file
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig fail: %v", err)
	}

	log.Println(env.AppName)

	if env.AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// construct DB connection
	dbConnection, err := models.RegisterDB(env)
	if err != nil {
		log.Fatalf("DB connection setup failed: %v", err)
	}

	router, err := router.RegisterRouter(env, dbConnection)
	if err != nil {
		log.Fatalf("Router setup failed: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", env.AppPort),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	} else {
		log.Println("Server exited properly")
	}
}
