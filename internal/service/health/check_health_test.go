package service

import (
	"context"
	"errors"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/repository/health/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		serviceName string
		instanceID  string

		setupRepoMock func(m *mocks.Health)

		expectNil      bool
		expectInstance string
		expectError    bool
	}{
		{
			name:        "repo error -> return nil",
			serviceName: "test-service",
			instanceID:  "instance-1",

			setupRepoMock: func(m *mocks.Health) {
				m.On("RedisPing", mock.Anything).
					Return(errors.New("redis down"))
			},

			expectNil:   true,
			expectError: true,
		},
		{
			name:        "repo healthy with predefined instanceID",
			serviceName: "test-service",
			instanceID:  "instance-1",

			setupRepoMock: func(m *mocks.Health) {
				m.On("RedisPing", mock.Anything).
					Return(nil)
			},

			expectNil:      false,
			expectInstance: "instance-1",
		},
		{
			name:        "repo healthy with empty instanceID -> generate new",
			serviceName: "test-service",
			instanceID:  "",

			setupRepoMock: func(m *mocks.Health) {
				m.On("RedisPing", mock.Anything).
					Return(nil)
			},

			expectNil:      false,
			expectInstance: "not-empty",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// ===== 1. Setup mock =====
			repoMock := mocks.NewHealth(t)
			tc.setupRepoMock(repoMock)

			// ===== 2. Init service =====
			s := &healthService{
				repo:        repoMock,
				serviceName: tc.serviceName,
				instanceId:  tc.instanceID,
			}

			// ===== 3. Call =====
			ctx := context.Background()
			result, err := s.CheckHealth(ctx)

			// ===== 4. Assert error =====
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)

			// ===== 5. Assert nil =====
			if tc.expectNil {
				assert.Nil(t, result)
				return
			}

			// ===== 6. Assert result =====
			assert.NotNil(t, result)
			assert.Equal(t, "OK", result.Message)
			assert.Equal(t, tc.serviceName, result.ServiceName)

			// ===== 7. Assert instance =====
			if tc.expectInstance == "not-empty" {
				assert.NotEmpty(t, result.InstanceID)
			} else {
				assert.Equal(t, tc.expectInstance, result.InstanceID)
			}
		})
	}
}
