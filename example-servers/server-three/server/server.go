package server

import (
	"net/http"
	"server-three/middleware"
	"server-three/server/routes"
)

func Start() {
	m := middleware.New("server-three")

	http.HandleFunc("/operation", m.Public(routes.Operation))

	http.ListenAndServe(":8082", nil)
}
