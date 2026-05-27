package service

import (
	"github.com/HOangAG2207/GoBeK03/internal/model"
	"github.com/HOangAG2207/GoBeK03/internal/utils"
)

func (s *healthService) CheckHealth() (*model.Health, error) {
	ok := s.repo.HealthPing()
	if !ok {
		return nil, nil
	}
	id := s.instanceId
	if id == "" {
		id = utils.UuidGenerator()
	}

	return &model.Health{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  id,
	}, nil
}
