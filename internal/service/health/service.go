package service

import (
	"github.com/HOangAG2207/GoBeK03/internal/model"
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
)

//go:generate mockery --name Health --filename check_health_mock.go --output ./mocks
type Health interface {
	CheckHealth() (*model.Health, error)
}

type healthService struct {
	repo        repository.Health
	serviceName string
	instanceId  string
}

func NewHealth(repo repository.Health, serviceName, instanceId string) Health {
	return &healthService{
		repo:        repo,
		serviceName: serviceName,
		instanceId:  instanceId,
	}
}
