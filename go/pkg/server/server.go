package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/remotehack/bottle/pkg/config"
	"github.com/remotehack/bottle/pkg/persister"
)

const (
	readTimeout     = 10
	shutdownTimeout = 5
)

type Server struct {
	router    *http.ServeMux
	config    config.Config
	persister persister.Persister
}

func New(cfg config.Config, persister persister.Persister) (Server, error) {
	return Server{
		router:    http.NewServeMux(),
		persister: persister,
		config:    cfg,
	}, nil
}

func (s *Server) Routes() {
	s.router.HandleFunc("/", s.writeRequest())
}

func (s *Server) Serve(ctx context.Context) {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", s.config.Port),
		Handler:           s.router,
		TLSConfig:         nil,
		ReadTimeout:       readTimeout * time.Second,
		ReadHeaderTimeout: readTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %s", err)
		}
	}()

	log.Println("server started")

	<-ctx.Done()

	log.Println("kill signal received")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed: %s", err)
	}

	log.Printf("server exited properly")
}

func (s *Server) writeRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.persister.Write("filename", "data")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "writing file failed: %s", err)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
