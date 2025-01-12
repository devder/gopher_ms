package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRouteExists(t *testing.T) {
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, chiRoutes chi.Router, route string) {
	found := false
	chi.Walk(chiRoutes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("route %s not found", route)
	}
}
