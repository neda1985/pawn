package main

import (
	"context"
	"dvb_pawn_shop/pkg/apis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shop := apis.NewShop(5)

		httpServer := &http.Server{
		Addr: ":8080" ,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	log.Print("http server started listening on port 8080" )
	http.HandleFunc("/offer", shop.HandleOffer)
	log.Print("Waiting request...")
	ctx, signalCancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer signalCancel()

	<-ctx.Done()

	log.Print("os.Interrupt - shutting down gracefully...\n")

	Context, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(Context); err != nil {
		log.Printf("shutdown error: %v\n", err)
		os.Exit(1)
	}
	log.Println("Context stopped")

}



