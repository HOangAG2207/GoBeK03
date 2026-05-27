package pkg_uuid

import "github.com/google/uuid"

//go:generate mockery --name Generator --filename generator_uuid_mock.go --output ./mocks
type Generator interface {
	Generate() string
}
type generator struct {
}

func NewGenerator() Generator {
	return &generator{}
}

func (g *generator) Generate() string {
	return uuid.New().String()
}
