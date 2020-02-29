package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MishaNiki/chat/internal/app/hub"
	"github.com/MishaNiki/chat/internal/app/model"
	"github.com/MishaNiki/chat/internal/app/templates"
	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	config    *Config
	logger    *logrus.Logger
	router    *http.ServeMux
	hub       *hub.Hub
	templates *templates.Templates
}

// New create new server
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: http.NewServeMux(),
		hub:    hub.New(),
	}
}

// Start starting server
func (s *Server) Start() error {

	var err error

	// server configuration

	s.configureRouter()

	if err = s.configureLogger(); err != nil {
		return err
	}

	if err = s.configureTemplates(); err != nil {
		return err
	}

	//server startup and safe shutdown
	s.logger.Info("Starting Server, port ", s.config.BindPort)

	go s.hub.Run()

	srv := http.Server{
		Addr:    s.config.BindPort,
		Handler: s.router,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		// setting a timeout to shut down the server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

	return err
}

/*
 *	SERVER CONFIGURATION SECTION
 */

// logger configuration
func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

// router configuration
func (s *Server) configureRouter() {
	// creation of routers
	stdAPIHandler := http.NewServeMux()
	//httprouterAPIHandler := httprouter.New()

	// Attaching url functions of handlers
	stdAPIHandler.HandleFunc("/", s.handleRoom())
	stdAPIHandler.HandleFunc("/ws", s.serveWS())

	// Creating a single router
	s.router.Handle("/", stdAPIHandler)

	// Create static storage
	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./web")),
	)
	s.router.Handle("/static/", staticHandler)
}

func (s *Server) configureTemplates() error {
	temp, err := templates.New(s.config.Templates)
	if err != nil {
		return err
	}
	s.templates = temp
	return nil
}

/*
 *	HANDLER FUNCTION SECTION
 */

func (s *Server) handleRoom() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		s.templates.Root.Execute(w, nil)
	}
}

/*
 / Work with web sockets
*/
func (s *Server) serveWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("conetect cliet")
		conn, err := s.hub.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Error("Error connection, s.hub.Upgrader.Upgrade")
			return
		}
		client := model.NewClient(conn)
		s.hub.Register <- client

		go s.hub.Client().WritePump(client)
		go s.hub.Client().ReadPump(client)
	}
}
