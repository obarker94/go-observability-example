package server

import (
	"net/http"
	"server-two/middleware"
	"server-two/server/routes"
)

func Start() {
	m := middleware.New("server-two")

	http.HandleFunc("/operation", m.Public(routes.Operation))

	http.ListenAndServe(":8081", nil)
}
