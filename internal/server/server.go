package server

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	http *http.Server
}

func New(endPoint string, handler http.Handler) *Server {
	var s Server
	s.http = &http.Server{
		Addr:    endPoint,
		Handler: handler,
	}
	return &s
}

func (s *Server) RunWithContext(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.http.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		timeOutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		return s.http.Shutdown(timeOutCtx)
	})
	return g.Wait()

}
