package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
	var req GetOwnerByUUIDRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	uuid := req.UUID

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
	var req GetOwnerByNameRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	name := req.Name

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
	owners, err := s.ownerRepo.GetAllOwners(ctx.Context(), 10, 0)
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
	var req DeleteOwnerRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadRequestError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if err := s.ownerRepo.DeleteOwner(ctx.Context(), req.UUID); err != nil {
		s.log.Error("Failed to delete owner", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   req.UUID,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
