package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"streaming-service/internal/service"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	apiGroup := app.Group("/v1")

	apiGroup.Post("/movies", r.Service.CreateMovie)
	apiGroup.Get("/movies/:id", r.Service.GetMovie)
	apiGroup.Get("/movies", r.Service.GetAllMovies)
	apiGroup.Put("/movies/:id", r.Service.UpdateMovie)
	apiGroup.Delete("/movies/:id", r.Service.DeleteMovie)

	return app
}
