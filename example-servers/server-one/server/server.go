package server

import (
	"net/http"
	"server-one/middleware"
	"server-one/server/routes"
)

func Start() {
	m := middleware.New("server-one")

	http.HandleFunc("/tracing", m.Public(routes.Example))

	http.ListenAndServe(":8080", nil)
}
