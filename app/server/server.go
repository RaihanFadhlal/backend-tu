package server

import (
	"fmt"
	"log"
	"net/http"
	"backendtku/app/routes"
	"backendtku/config"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	config *config.Config
	router *mux.Router
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		router: mux.NewRouter(),
	}
}

func (s *Server) Run() {
	routes.InitializeRoutes(s.router, s.config.DB)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(s.router)
	fmt.Printf("Listening to port %s\n", s.config.AppPort)
	log.Fatal(http.ListenAndServe(":"+s.config.AppPort, handler))
}
