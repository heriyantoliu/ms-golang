package service

import "net/http"

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"HealthCheck",
		"Get",
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.Write([]byte("OK"))
		},
	},
}