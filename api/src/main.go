package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow only localhost:3000
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	// isProd     = flag.Bool("prod", false, "is production mode")
	httpPort   = flag.Int("httpPort", 4000, "http port")
	httpServer *http.Server
)

type application struct {
	logger *slog.Logger
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v: %v", err, string(debug.Stack()))
	} else {
		log.Println("Auth server has been gracefully shutdown")
	}
}

func run() error {
	flag.Parse()

	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	app := &application{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", *httpPort),
		Handler:      enableCORS(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error)

	// start http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()
	log.Printf("HTTP server is running at :%d", *httpPort)

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM,
	)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		log.Printf("signal '%v' detected, server is being shutdown", sig)
		cancel()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	}
}

// func newAuthDbClient(ctx context.Context) (*dynamodb.Client, error) {
// 	awsConfig, err := config.LoadDefaultConfig(ctx)
// 	if err != nil {
// 		return &dynamodb.Client{}, err
// 	}

// 	localAuthDbEndpoint := "http://authdb:8000"
// 	// wait for authdb instance to spin up. 'depends_on' attribute in docker compose file doesn't work
// 	if !*isProd {
// 		client := &http.Client{}
// 		for i := 0; i < 10; i++ {
// 			resp, err := client.Get(localAuthDbEndpoint)
// 			if err != nil {
// 				if i == 9 {
// 					return &dynamodb.Client{}, fmt.Errorf("cannot connect to local authdb: %v", err)
// 				}
// 				time.Sleep(time.Second)
// 				continue
// 			}
// 			resp.Body.Close()
// 			break
// 		}
// 	}

// 	c := dynamodb.NewFromConfig(awsConfig,
// 		func(o *dynamodb.Options) {
// 			if !*isProd {
// 				o.BaseEndpoint = aws.String(localAuthDbEndpoint)
// 			}
// 		},
// 	)

// 	return c, nil
// }
