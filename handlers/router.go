package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func CreateRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CRSF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(router chi.Router) {

		// version 1
		router.Route("/v1", func(router chi.Router) {

			router.Get("/healthcheck", healthCheck)
			router.Get("/todos", getTodos)
			router.Get("/todos/{id}", getTodoById)
			router.Get("/todos/org", getTodoByOrg) // Filter by Organisation Name
			router.Get("/todos/vol", getTodoByVol) // Filter by Volunteer Type
			router.Post("/todos/create", createTodo)
			router.Put("/todos/update/{id}", updateTodo)
			router.Delete("/todos/delete/{id}", deleteTodo)
			router.Post("/signup", SignupHandler)
			router.Post("/login", LoginHandler)
			router.Get("/users", GetUserByIDHandler)

		})

		// version 2 - add it if you want
		// router.Route("/v2", func(router chi.Router) {
		// })

	})

	return router

}
