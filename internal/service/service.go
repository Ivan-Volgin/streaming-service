package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"streaming-service/internal/dto"
	"streaming-service/internal/repo"
)

type Service interface {
	CreateMovie(c *fiber.Ctx) error
	GetMovie(c *fiber.Ctx) error
	GetAllMovies(c *fiber.Ctx) error
	UpdateMovie(c *fiber.Ctx) error
	DeleteMovie(c *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
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
	movieID, err := s.repo.CreateMovie(ctx.Context(), &movie)
	if err != nil {
		s.log.Error("Failed to create movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   map[string]interface{}{"movieID": movieID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetMovie(ctx *fiber.Ctx) error {
	//var req GetMovieRequest
	//
	//if err := json.Unmarshal(c.Body(), &req); err != nil {
	//	s.log.Error("Invalid request body", zap.Error(err))
	//	return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	//}
	//
	//uuid := req.UUID
	//
	//movie, err := s.repo.GetMovieByID(ctx.Context(), uuid)
	return nil
}

func (s *service) GetAllMovies(c *fiber.Ctx) error {
	return nil
}

func (s *service) UpdateMovie(c *fiber.Ctx) error {
	return nil
}

func (s *service) DeleteMovie(c *fiber.Ctx) error {
	return nil
}
