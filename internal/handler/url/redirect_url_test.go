package url

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/service/url/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_RedirectURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx echo.Context)
		setupMockService func(ctx context.Context) *mocks.UrlService

		expectedStatus int
		expectedUrl    string
	}{
		{
			name: "normal case - redirect successfully",
			setupRequest: func(ctx echo.Context) {
				req := httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/abc123",
					nil,
				)
				ctx.SetRequest(req)
				ctx.SetParamNames("code")
				ctx.SetParamValues("abc123")
			},
			setupMockService: func(ctx context.Context) *mocks.UrlService {
				serviceMock := mocks.NewUrlService(t)
				serviceMock.On("GetURL", ctx, "abc123").
					Return("https://google.com", nil).
					Once()
				return serviceMock
			},
			expectedStatus: http.StatusFound,
			expectedUrl:    "https://google.com",
		},
		{
			name: "missing code parameter",

			setupRequest: func(ctx echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				ctx.SetRequest(req)
				// No code parameter set
			},

			setupMockService: func(ctx context.Context) *mocks.UrlService {
				return mocks.NewUrlService(t)
			},

			expectedStatus: http.StatusBadRequest,
			expectedUrl:    "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Add test logic here
			e := echo.New()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := e.NewContext(req, rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx.Request().Context())

			testHandler := NewUrlHandler(mockSvc)

			testHandler.RedirectURL(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedUrl, rec.Header().Get("Location"))
		})
	}
}
