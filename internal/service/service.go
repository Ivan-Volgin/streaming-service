package service

import (
	"go.uber.org/zap"
	"streaming-service/internal/repo"
)

type service struct {
	movieRepo repo.MovieRepository
	ownerRepo repo.OwnerRepository
	log       *zap.SugaredLogger
}

type Service interface {
	MovieService
	OwnerService
}

func NewService(movieRepo repo.MovieRepository, ownerRepo repo.OwnerRepository, logger *zap.SugaredLogger) Service {
	return &service{
		movieRepo: movieRepo,
		ownerRepo: ownerRepo,
		log:       logger,
	}
}
