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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tasjen/fz/db"
)

var (
	// isProd     = flag.Bool("prod", false, "is production mode")
	httpPort    = flag.String("httpPort", os.Getenv("PORT"), "http port")
	postgresURI = flag.String("postgresURI", os.Getenv("POSTGRESQL_URI"), "PostgresQL URI")
	httpServer  *http.Server
)

type application struct {
	logger *slog.Logger
	db     *db.Queries
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v: %v", err, string(debug.Stack()))
	} else {
		log.Println("Server has been gracefully shutdown")
	}
}

func run() error {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	app := &application{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	pool, err := newPostgresClient(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()
	app.db = db.New(pool)

	// queries := db.New(pool)

	// user, err := queries.GetUser(ctx, pgtype.UUID{Bytes: uuid.New(), Valid: true})
	// if err != nil {
	// 	return fmt.Errorf("failed to get user: %w", err)
	// }
	// log.Printf("User: %+v", user)

	httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", *httpPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return startHttpServer()
}

func startHttpServer() error {
	errCh := make(chan error)
	// start http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()
	log.Printf("HTTP server is running at :%s", *httpPort)

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
