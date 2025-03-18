package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"streaming-service/internal/dto"
	"streaming-service/internal/repo"
)

type OwnerService interface {
	CreateOwner(ctx *fiber.Ctx) error
	GetOwnerByUUID(ctx *fiber.Ctx) error
	GetOwnerByName(ctx *fiber.Ctx) error
	GetAllOwners(ctx *fiber.Ctx) error
	UpdateOwner(ctx *fiber.Ctx) error
	DeleteOwner(ctx *fiber.Ctx) error
}

func (s *service) CreateOwner(ctx *fiber.Ctx) error {
	var req CreateOwnerRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	owner := repo.Owner{
		Name: req.Name,
	}
	ownerID, err := s.ownerRepo.CreateOwner(ctx.Context(), &owner)
	if err != nil {
		s.log.Error("Failed to create owner", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   map[string]string{"movieID": ownerID},
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (s *service) GetOwnerByUUID(ctx *fiber.Ctx) error {
	uuid := ctx.Params("uuid")
	if uuid == "" {
		s.log.Error("Missing UUID in URL parameters", zap.String("uuid", uuid))
		return dto.BadRequestError(ctx, dto.FieldRequired, "Missing UUID in URL parameters")
	}

	owner, err := s.ownerRepo.GetOwnerByID(ctx.Context(), uuid)
	if err != nil {
		s.log.Error("Failed to get owner", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: map[string]interface{}{
			"uuid":       owner.UUID,
			"name":       owner.Name,
			"created_at": owner.Created_at,
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetOwnerByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if name == "" {
		s.log.Error("Missing name in URL parameters", zap.String("name", name))
		return dto.BadRequestError(ctx, dto.FieldRequired, "Missing name in URL parameters")
	}

	owner, err := s.ownerRepo.GetOwnerByName(ctx.Context(), name)
	if err != nil {
		s.log.Error("Failed to get owner", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: map[string]interface{}{
			"uuid":       owner.UUID,
			"name":       owner.Name,
			"created_at": owner.Created_at,
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetAllOwners(ctx *fiber.Ctx) error {
	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		s.log.Error("Invalid limit parameter", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid or missing 'limit' parameter")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		s.log.Error("Invalid offset parameter", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid or missing 'offset' parameter")
	}
	owners, err := s.ownerRepo.GetAllOwners(ctx.Context(), limit, offset)
	if err != nil {
		s.log.Error("Failed to get owners", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	response := dto.Response{
		Status: "success",
		Data:   owners,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateOwner(ctx *fiber.Ctx) error {
	var req UpdateOwnerRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	updatedOwner := repo.Owner{
		Name: req.Name,
	}

	if err := s.ownerRepo.UpdateOwner(ctx.Context(), req.UUID, &updatedOwner); err != nil {
		s.log.Error("Failed to update movie", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: map[string]interface{}{
			"ownerUUID": req.UUID,
			"owner":     updatedOwner,
		},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteOwner(ctx *fiber.Ctx) error {
	uuid := ctx.Params("uuid")
	if uuid == "" {
		s.log.Error("Missing UUID in URL parameters", zap.String("uuid", uuid))
		return dto.BadRequestError(ctx, dto.FieldRequired, "Missing UUID in URL parameters")
	}

	if err := s.ownerRepo.DeleteOwner(ctx.Context(), uuid); err != nil {
		s.log.Error("Failed to delete owner", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   uuid,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
