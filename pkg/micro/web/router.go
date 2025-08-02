package web

import (
	"github.com/justinas/alice"
	"net/http"
)

// Route -- Defines a single route, e.g. a human readable name, HTTP method,
// pattern the function that will execute when the route is called.
type Route struct {
	Name        string              `json:"-"`
	Method      string              `json:"m"`
	Pattern     string              `json:"p"`
	Handler     http.Handler        `json:"-"`
	Middlewares []alice.Constructor `json:"-"`
	AuthInfo    AuthInfo            `json:"a"`
	Timeout     int64               `json:"t"`
}

// AuthInfo -- authentication and authorization for route
type AuthInfo struct {
	Enable             bool            `json:"e"`
	TokenType          string          `json:"t"`
	RequirePermissions map[int64]int64 `json:"r"`
}

// Routes -- Defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route
