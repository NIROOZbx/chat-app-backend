package main

import (
	"chat-app/internal/bootstrap"
	"chat-app/internal/shared/config"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Println("Config loaded ✅")

	app := bootstrap.StartUp(cfg)

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	 go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    
 

	//graceful shutdown

	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, signals...)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

}
