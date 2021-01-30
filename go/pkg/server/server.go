package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/remotehack/bottle/pkg/config"
	"github.com/remotehack/bottle/pkg/persister"
	"github.com/remotehack/bottle/pkg/serializer"
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
		fn, err := s.getFilename(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "server.writeRequest: %s", err)

			return
		}

		err = s.persister.Write(fn, serializer.Serialize(fmt.Sprintf("%#v", r.Header)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "writing file failed: %s", err)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) getFilename(r *http.Request) (string, error) {
	fmt.Printf("request url: %s\nconfigured url: %s\n\n", r.URL.Host, s.config.Host)
	fmt.Printf("whole %#v\n", r.URL)

	if strings.HasSuffix(r.URL.Host, s.config.Host) {
		h := strings.TrimRight(r.URL.Host, s.config.Host)
		if h == "" {
			return "", errors.New("server.getFileName there was no subdomain")
		}

		return h, nil
	}

	return "development", nil
}
