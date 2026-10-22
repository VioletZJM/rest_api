package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/rest_api/database"
	"example.com/rest_api/routes"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.LoadJWTSecret()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	routes.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("server shutdown failed: %v", err)
		srv.Close()
	}

	log.Println("Server exited")

}
