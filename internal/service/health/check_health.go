package service

import (
	"github.com/HOangAG2207/GoBeK03/internal/model"
)

func (s *healthService) CheckHealth() (*model.Health, error) {
	ok := s.repo.HealthPing()
	if !ok {
		return nil, nil
	}
	id := s.instanceId
	if id == "" {
		id = s.uuidGen.Generate()
	}

	return &model.Health{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  id,
	}, nil
}
