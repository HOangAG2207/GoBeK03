package repository

//go:generate mockery --name Health --filename check_health_mock.go --output ./mocks
type Health interface {
	HealthPing() bool
}

type healthRepository struct {
}

func NewHealth() Health {
	return &healthRepository{}
}
