package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"streaming-service/internal/service"
)

type Routers struct {
	MovieService service.MovieService
	OwnerService service.OwnerService
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

	apiGroup.Post("/movies", r.MovieService.CreateMovie)
	apiGroup.Get("/movies/:id", r.MovieService.GetMovie)
	apiGroup.Get("/movies", r.MovieService.GetAllMovies)
	apiGroup.Put("/movies/", r.MovieService.UpdateMovie)
	apiGroup.Delete("/movies/:id", r.MovieService.DeleteMovie)

	apiGroup.Post("/owners", r.OwnerService.CreateOwner)
	apiGroup.Get("/owners/id/:id", r.OwnerService.GetOwnerByUUID)
	apiGroup.Get("/owners/name/:name", r.OwnerService.GetOwnerByName)
	apiGroup.Get("/owners", r.OwnerService.GetAllOwners)
	apiGroup.Put("/owners/", r.OwnerService.UpdateOwner)
	apiGroup.Delete("/owners/:id", r.OwnerService.DeleteOwner)

	return app
}
