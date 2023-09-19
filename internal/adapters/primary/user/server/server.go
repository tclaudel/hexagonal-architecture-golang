//go:generate oapi-codegen -generate types -package server -o ./types.gen.go ../../../../../api/users.swagger.yaml
//go:generate oapi-codegen -generate chi-server -package server -o ./server.gen.go ../../../../../api/users.swagger.yaml
//go:generate oapi-codegen -generate client -package server -o ./client.gen.go ../../../../../api/users.swagger.yaml

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

var (
	ErrInvalidServerPort       = errors.New("invalid server port")
	ErrListenAndServeAppServer = errors.New("error while listening and serving app server")
	ErrShutDownAppServer       = errors.New("error while shutting down app server")
)

type Server struct {
	server *http.Server
}

type Params struct {
	UserUseCase primary.UserUseCase
	Port        string
}

func NewServer(_ context.Context, params Params) (*Server, error) {
	address, err := netip.ParseAddrPort(fmt.Sprintf("0.0.0.0:%s", params.Port))
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrInvalidServerPort, err)
	}

	server := Server{server: &http.Server{
		Addr: address.String(),
		Handler: HandlerWithOptions(NewHandlers(params.UserUseCase), ChiServerOptions{
			BaseURL:          "",
			BaseRouter:       chi.NewRouter(),
			Middlewares:      nil,
			ErrorHandlerFunc: ErrorHandlerFunc,
		}),
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}}

	return &server, nil
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	log.Printf("starting app server %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%s, %s", ErrListenAndServeAppServer, err)

		return fmt.Errorf("%w, %s", ErrListenAndServeAppServer, err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Printf("error while shutting down app server, %s", err)

		return fmt.Errorf("%w, %s", ErrShutDownAppServer, err)
	}

	return nil
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, _ *http.Request, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(buf.Bytes()) //nolint:errcheck
}

// ErrorHandlerFunc is used by the generated code to return json errors.
func ErrorHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	JSON(w, r, http.StatusBadRequest, Error{Message: err.Error()})
}
