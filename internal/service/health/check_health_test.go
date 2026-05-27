package service

import (
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/repository/health/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serviceName    string
		instanceID     string
		setupRepoMock  func(m *mocks.Health)
		expectNil      bool
		expectInstance string
	}{
		{
			name:        "repo not healthy -> return nil",
			serviceName: "test-service",
			instanceID:  "instance-1",
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(false)
			},
			expectNil: true,
		},
		{
			name:        "repo healthy with predefined instanceID",
			serviceName: "test-service",
			instanceID:  "instance-1",
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(true)
			},
			expectNil:      false,
			expectInstance: "instance-1",
		},
		{
			name:        "repo healthy with empty instanceID -> generate uuid",
			serviceName: "test-service",
			instanceID:  "",
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(true)
			},
			expectNil:      false,
			expectInstance: "not-empty", // chỉ check != ""
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoMock := mocks.NewHealth(t)
			tc.setupRepoMock(repoMock)

			s := &healthService{
				repo:        repoMock,
				serviceName: tc.serviceName,
				instanceId:  tc.instanceID,
			}

			result, err := s.CheckHealth()

			assert.NoError(t, err)

			if tc.expectNil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, "OK", result.Message)
			assert.Equal(t, tc.serviceName, result.ServiceName)

			if tc.expectInstance == "not-empty" {
				assert.NotEmpty(t, result.InstanceID)
			} else {
				assert.Equal(t, tc.expectInstance, result.InstanceID)
			}
		})
	}
}
