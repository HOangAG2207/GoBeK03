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

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockRepo.UrlRepository

		inputOriginalURL string

		expectedCodeLen int
		expectedError   error
	}{
		{
			name: "success_first_try",
			setupMockRepo: func(ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				repo.
					On("StoreURLIfAbsent", ctx, mock.AnythingOfType("string"), "https://example.com").
					Return(true, nil).
					Once()

				return repo
			},
			inputOriginalURL: "https://example.com",
			expectedCodeLen:  10,
			expectedError:    nil,
		},
		{
			name: "retry_then_success",
			setupMockRepo: func(ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				// fail lần 1
				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com").
					Return(false, nil).
					Once()

				// success lần 2
				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com").
					Return(true, nil).
					Once()

				return repo
			},
			inputOriginalURL: "https://example.com",
			expectedCodeLen:  10,
			expectedError:    nil,
		},
		{
			name: "repo_error",
			setupMockRepo: func(ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				repo.
					On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com").
					Return(false, errors.New("db error")).
					Once()

				return repo
			},
			inputOriginalURL: "https://example.com",
			expectedCodeLen:  0,
			expectedError:    errors.New("db error"),
		},
		{
			name: "max_retry_exceeded",
			setupMockRepo: func(ctx context.Context) *mockRepo.UrlRepository {
				repo := mockRepo.NewUrlRepository(t)

				for i := 0; i < maxRetryAttempts; i++ {
					repo.
						On("StoreURLIfAbsent", ctx, mock.Anything, "https://example.com").
						Return(false, nil).
						Once()
				}

				return repo
			},
			inputOriginalURL: "https://example.com",
			expectedCodeLen:  0,
			expectedError:    ErrMaxRetriesExceeded,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo := tc.setupMockRepo(ctx)

			service := NewUrlService(repo, 10)

			code, err := service.ShortenURL(ctx, tc.inputOriginalURL)

			// check error
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)

			// check code length (do random)
			assert.Len(t, code, tc.expectedCodeLen)
		})
	}
}
