package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"

	"streaming-service/internal/dto"
	"streaming-service/internal/repo"
)

type MovieService interface {
	CreateMovie(c *fiber.Ctx) error
	GetMovie(c *fiber.Ctx) error
	GetAllMovies(c *fiber.Ctx) error
	UpdateMovie(c *fiber.Ctx) error
	DeleteMovie(c *fiber.Ctx) error
}

func (s *service) CreateMovie(ctx *fiber.Ctx) error {
	var req CreateMovieRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	movie := repo.Movie{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Year:        req.Year,
	}
	movieID, err := s.movieRepo.CreateMovie(ctx.Context(), &movie, req.OwnerName)
	if err != nil {
		s.log.Error("Failed to create movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   map[string]string{"movieID": movieID},
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (s *service) GetMovie(ctx *fiber.Ctx) error {
	var req GetMovieRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	uuid := req.UUID

	movie, err := s.movieRepo.GetMovieByID(ctx.Context(), uuid)
	if err != nil {
		s.log.Error("Failed to get movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: map[string]string{
			"title":       movie.Title,
			"author":      movie.Author,
			"description": movie.Description,
			"year":        strconv.Itoa(movie.Year),
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetAllMovies(ctx *fiber.Ctx) error {
	movies, err := s.movieRepo.GetAllMovies(ctx.Context(), 10, 0)
	if err != nil {
		s.log.Error("Failed to get movies", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	response := dto.Response{
		Status: "success",
		Data:   movies,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateMovie(ctx *fiber.Ctx) error {
	var req UpdateMovieRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	updatedMovie := repo.Movie{
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		Year:        req.Year,
	}

	if err := s.movieRepo.UpdateMovie(ctx.Context(), req.UUID, &updatedMovie); err != nil {
		s.log.Error("Failed to update movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: map[string]interface{}{
			"movieUUID": req.UUID,
			"movie":     updatedMovie,
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteMovie(ctx *fiber.Ctx) error {
	var req DeleteMovieRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if err := s.movieRepo.DeleteMovie(ctx.Context(), req.UUID); err != nil {
		s.log.Error("Failed to delete movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   req.UUID,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
