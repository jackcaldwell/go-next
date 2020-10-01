package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackcaldwell/go-next/pkg/auth"
	"github.com/jackcaldwell/go-next/pkg/session"
	"github.com/jackcaldwell/go-next/pkg/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Setup config flags
	port := flag.String("port", "8080", "The port that the web server should run on")
	dbPort := flag.String("dbport", "5432", "The port that the application database is running on")
	dbHost := flag.String("dbhost", "localhost", "The host that the application database is running on")
	dbUser := flag.String("dbuser", "postgres", "The username for database")
	dbPass := flag.String("dbpass", "postgres", "The password for the database")
	dbName := flag.String("dbname", "gonext", "The database name")
	dbSSLMode := flag.String("dbsslmode", "disable", "The SSL mode that the database connection should use")
	flag.Parse()

	// Apply database migrations
	m, err := migrate.New(
		"file://../db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", *dbUser, *dbPass, *dbHost, *dbPort, *dbName, *dbSSLMode),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		panic(err)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", *dbUser, *dbPass, *dbName, *dbHost, *dbPort, *dbSSLMode))

	if err != nil {
		panic(err)
	}

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	httpLogger := log.With(logger, "component", "http")
	r := mux.NewRouter()

	var (
		sessions = session.NewRepository(db)
		users    = user.NewRepository(db)
	)

	var authService auth.Service
	authService = auth.NewService(sessions, users)
	auth.RegisterRoutes(authService, httpLogger, r)

	http.Handle("/", accessControl(r))

	errs := make(chan error, 2)

	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	go func() {
		logger.Log("transport", "http", "address", ":"+*port, "msg", "listening")
		errs <- http.ListenAndServe(":"+*port, handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r))
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func removeURLTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
