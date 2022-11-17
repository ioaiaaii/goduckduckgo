/*
https://medium.com/rungo/creating-a-simple-hello-world-http-server-in-go-31c7fd70466e
https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f
*/

package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	mux  *http.ServeMux
	srv  *http.Server
	opts options
}

// Returns a pointer to the newly created Server
func NewServer(opts ...Option) *Server {

	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	httpHandler := http.NewServeMux()

	return &Server{
		mux: httpHandler,
		srv: &http.Server{
			Addr:         ":" + options.HTTPPort,
			ReadTimeout:  options.HTTPServerTimeout * time.Second,
			WriteTimeout: options.HTTPServerTimeout * time.Second,
			IdleTimeout:  2 * options.HTTPServerTimeout * time.Second,
			Handler:      httpHandler,
		},
		opts: options,
	}

}

// Handle registers the handler for the given pattern.
func (s *Server) HandleAPIPath(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) ListenAndServe() error {

	log.Printf("listening for serving HTTP, address: %v", s.srv.Addr)

	// Run server
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
		return err
	}

	return nil
}

func (s *Server) Shutdown(err error) {

	if err == http.ErrServerClosed {
		log.Println("msg", "internal server closed unexpectedly")
		return
	}

	timeOut := s.opts.HTTPServerShutdownTimeout * time.Second
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), timeOut)
	defer cancelShutdown()

	if err := s.srv.Shutdown(gracefulCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
	} else {
		log.Printf("gracefully stopped\n")
	}

	log.Printf("shutdown gracefully...")
}
