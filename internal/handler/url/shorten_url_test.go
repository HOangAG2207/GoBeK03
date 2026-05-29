package url

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/service/url/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		requestBody string

		setupMockSvc func() *mocks.UrlService

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			requestBody: `{"url":"https://google.com"}`,

			setupMockSvc: func() *mocks.UrlService {
				mockSvc := &mocks.UrlService{}
				mockSvc.
					On("ShortenURL", mock.Anything, "https://google.com").
					Return("abc123", nil)
				return mockSvc
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abc123","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "invalid payload",

			requestBody: `{}`,

			setupMockSvc: func() *mocks.UrlService {
				return &mocks.UrlService{}
			},

			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"code":"","message":"invalid request payload"}`,
		},
		{
			name: "service error",

			requestBody: `{"url":"https://google.com"}`,

			setupMockSvc: func() *mocks.UrlService {
				mockSvc := &mocks.UrlService{}
				mockSvc.
					On("ShortenURL", mock.Anything, "https://google.com").
					Return("", errors.New("redis down"))
				return mockSvc
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"code":"","message":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/api/url/shorten", bytes.NewBufferString(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			mockSvc := tc.setupMockSvc()
			handler := NewUrlHandler(mockSvc)

			err := handler.ShortenURL(ctx)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}

}
