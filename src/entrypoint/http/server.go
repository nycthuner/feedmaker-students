package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"gorm.io/gorm"
)

type Server struct {
	app        *gin.Engine
	httpServer *http.Server
	db         *gorm.DB
}

func NewServer(port string, db *gorm.DB) *Server {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	return &Server{
		app,
		&http.Server{
			Addr: ":" + port,
		},
		db,
	}
}
func (server *Server) Start() error {
	server.corsConfig()
	server.registerRoutes(server.db)

	fmt.Println("server started")

	server.initialize()
	if err := server.shutdown(); err != nil {
		return err
	}

	return nil
}

func (server *Server) corsConfig() {
	headers := handlers.AllowedHeaders([]string{
		"Accept",
		"Access-Control-Allow-Origin",
		"Accept-Encoding",
		"Accept-Language",
		"Authorization",
		"Connection",
		"Content-Length",
		"Content-Type",
		"Origin",
	})

	origins := handlers.AllowedOrigins([]string{
		"http://localhost:5173",
	})

	methods := handlers.AllowedMethods([]string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	})

	credentials := handlers.AllowCredentials()

	corsHandler := handlers.CORS(
		headers,
		origins,
		methods,
		credentials,
	)(server.app)

	server.httpServer.Handler = corsHandler
}

func (server *Server) registerRoutes(db *gorm.DB) {
	NewRoutes(
		server.app,
		db,
	).Register()
}

func (server *Server) initialize() {
	go func() {
		err := server.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (server *Server) shutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	fmt.Println("kill signal received, shutting down server...")

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(context); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	fmt.Println("server has been successfully shut down")

	return nil
}
