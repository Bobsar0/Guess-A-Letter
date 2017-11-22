package apiapp

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Server represents an HTTP server.
type Server struct {
	ln net.Listener

	// Handler to serve.
	Handler *Handler

	// Bind address to open.
	Addr string
}

func NewServer(address string, h *Handler) *Server {
	return &Server{
		Handler: h,
		Addr:    ":" + address,
	}
}

// Open opens a socket and serves the HTTP server.
//,
func (s *Server) Open(done chan bool, sigs chan os.Signal) error {
	//_, _ = done, sigs

	//[to make an app engine app uncomment the next two lines]
	// http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, s.Handler))
	// appengine.Main()

	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln
	// Start HTTP server.
	go func() {

		log.Fatal(http.Serve(s.ln, handlers.CombinedLoggingHandler(os.Stderr, s.Handler)))
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	// http.Serve(s.ln, handlers.CombinedLoggingHandler(os.Stderr, s.Handler))

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}
