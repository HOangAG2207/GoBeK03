package service

import (
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
	pkg_uuid "github.com/HOangAG2207/GoBeK03/pkg/uuid"
)

type Health interface {
}

type healthService struct {
	repo        repository.Health
	serviceName string
	instanceId  string
	uuidGen     pkg_uuid.Generator
}

func NewHealth(repo repository.Health, serviceName, instanceId string, uuidGen pkg_uuid.Generator) Health {
	return &healthService{
		repo:        repo,
		serviceName: serviceName,
		instanceId:  instanceId,
		uuidGen:     uuidGen,
	}
}
