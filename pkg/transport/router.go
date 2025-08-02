package transhttp

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/micro/web"
	middleware "github.com/nhdms/base-go/pkg/middlewares"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 10000 // ms
	NoTimeout      = -1
)

func InitRoutes(svc web.Service, routes web.Routes, path string) {
	globalTimeout := config.ViperGetInt64WithDefault("api.timeout", DefaultTimeout)
	router := mux.NewRouter()
	for _, route := range routes {
		//fullPath := basePath + route.Pattern
		// Start with the handler
		chain := alice.New(middleware.NewRequestID("api"))
		// Add timeout middleware if set
		if route.Timeout != NoTimeout {
			timeout := globalTimeout
			if route.Timeout > 0 {
				timeout = route.Timeout
			}
			chain = chain.Append(TimeoutMiddleware(time.Duration(timeout) * time.Millisecond))
		}

		// Append additional middlewares from the route definition
		chain = chain.Append(route.Middlewares...)
		// Return the final handler with all middlewares applied
		handler := chain.Then(route.Handler)
		router.Handle(route.Pattern, handler).Methods(route.Method)
	}
	svc.Handle("/health-check", new(HealthCheckHandler))
	svc.Handle("/", router)
}

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RespondMessage(w, http.StatusOK, "i'm ok!")
}

// TimeoutMiddleware adds a timeout to the route handler
func TimeoutMiddleware(duration time.Duration) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			// Replace request context with the timeout context
			r = r.WithContext(ctx)

			// Run the handler
			next.ServeHTTP(w, r)
		})
	}
}
