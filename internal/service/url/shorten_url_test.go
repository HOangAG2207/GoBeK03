package url

import (
	"context"
	"errors"
	"testing"

	mockRepo "github.com/HOangAG2207/GoBeK03/internal/repository/url/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_ShortenURL(t *testing.T) {
	t.Parallel()

	var ErrDB = errors.New("db error")

	testCases := []struct {
		name string

		inputURL string

		setupMock func(t *testing.T, ctx context.Context) *mockRepo.UrlRepository

		expectErr error
		expectLen int
	}{
		{
			name:     "success_first_try",
			inputURL: "https://example.com",

			setupMock: func(t *testing.T, ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com", int64(5)).
					Return(true, nil).
					Once()

				return repo
			},

			expectErr: nil,
			expectLen: 10,
		},
		{
			name:     "retry_then_success",
			inputURL: "https://example.com",

			setupMock: func(t *testing.T, ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com", int64(5)).
					Return(false, nil).
					Once()

				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com", int64(5)).
					Return(true, nil).
					Once()

				return repo
			},

			expectErr: nil,
			expectLen: 10,
		},
		{
			name:     "repo_error",
			inputURL: "https://example.com",

			setupMock: func(t *testing.T, ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com", int64(5)).
					Return(false, ErrDB).
					Once()

				return repo
			},

			expectErr: ErrDB,
			expectLen: 0,
		},
		{
			name:     "max_retry_exceeded",
			inputURL: "https://example.com",

			setupMock: func(t *testing.T, ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				for i := 0; i < maxRetryAttempts; i++ {
					repo.
						On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com", int64(5)).
						Return(false, nil).
						Once()
				}

				return repo
			},

			expectErr: ErrMaxRetriesExceeded,
			expectLen: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo := tc.setupMock(t, ctx)
			service := NewUrlService(repo, 10)

			code, err := service.ShortenURL(ctx, tc.inputURL, 5)

			// ===== assert error =====
			if tc.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectErr) // ✅ giờ sẽ PASS
				assert.Empty(t, code)
				return
			}

			assert.NoError(t, err)

			// ===== assert success =====
			assert.Len(t, code, tc.expectLen)
		})
	}
}
