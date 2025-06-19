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

// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/handler/extension"
// 	"github.com/99designs/gqlgen/graphql/handler/lru"
// 	"github.com/99designs/gqlgen/graphql/handler/transport"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	"github.com/tasjen/fz/graph"
// 	"github.com/vektah/gqlparser/v2/ast"
// )

// const defaultPort = "8080"

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

// 	srv.AddTransport(transport.Options{})
// 	srv.AddTransport(transport.GET{})
// 	srv.AddTransport(transport.POST{})

// 	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

// 	srv.Use(extension.Introspection{})
// 	srv.Use(extension.AutomaticPersistedQuery{
// 		Cache: lru.New[string](100),
// 	})

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }
