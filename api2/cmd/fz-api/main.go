package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tasjen/sathu/api/db"
)

var (
	// isProd     = flag.Bool("prod", false, "is production mode")
	httpPort    = flag.String("httpPort", os.Getenv("PORT"), "http port")
	postgresURI = flag.String("postgresURI", os.Getenv("POSTGRESQL_URI"), "PostgresQL URI")
)

// type application struct {
// 	logger *slog.Logger
// }

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v: %v", err, string(debug.Stack()))
		os.Exit(1)
	} else {
		log.Println("Server has been gracefully shutdown")
	}
}

func run() error {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := newPostgresClient(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := db.NewUserRepository(pool)
	userService := db.NewUserService(userRepo)
	userHandler := db.NewUserHandler(userService)
	router, err := db.NewRouter(&db.Config{}, userHandler)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	// start http server
	go func() {
		log.Printf("HTTP server is running at :%s\n", router.Addr)
		if err := router.Serve(); err != nil {
			errCh <- err
		}
	}()

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
		if err := httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	}
}

func startHttpServer(db *db.Queries) error {
	handlers := handlers.NewHandlers(db)
	routes := routes.NewRouter(handlers)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", *httpPort),
		Handler:      routes.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error)
	// start http server
	go func() {
		log.Printf("HTTP server is running at :%s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

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
		if err := httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	}
}

func newPostgresClient(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, *postgresURI)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
