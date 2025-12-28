package api

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"

	"github.com/mickamy/contest/example/internal/domain/session"
)

type Router = func(mux *http.ServeMux, options ...connect.HandlerOption)

func NewServer() http.Server {
	api := http.NewServeMux()
	for _, route := range []Router{
		session.Route,
	} {
		route(api)
	}

	mux := http.NewServeMux()
	corsHandler := cors.AllowAll().Handler(mux)

	return http.Server{
		Addr:    fmt.Sprintf(":%d", 8008),
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),
	}
}
