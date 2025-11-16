package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	client "github.com/fingertips18/fingertips18.github.io/backend/internal/client/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type imageRepositoryTestFixture struct {
	t               *testing.T
	mockHttpAPI     *client.MockHttpAPI
	imageRepository imageRepository
}

func newImageRepositoryTestFixture(t *testing.T) *imageRepositoryTestFixture {
	mockHttpAPI := client.NewMockHttpAPI(t)
	imageRepository := &imageRepository{
		uploadthingToken: "test_token_xxx",
		httpAPI:          mockHttpAPI,
	}

	return &imageRepositoryTestFixture{
		t:               t,
		mockHttpAPI:     mockHttpAPI,
		imageRepository: *imageRepository,
	}
}

func TestImageRepository_Upload(t *testing.T) {
	httpErr := errors.New("http error")
	missingFilesErr := errors.New("failed to validate image: files missing")
	missingNameErr := errors.New("failed to validate image: file[0]: name missing")
	invalidSizeErr := errors.New("failed to validate image: file[0]: size invalid")
	missingTypeErr := errors.New("failed to validate image: file[0]: type missing")
	invalidACLErr := errors.New("failed to validate image: acl must be 'public-read' or 'private'")
	invalidContentDispositionErr := errors.New("failed to validate image: contentDisposition must be 'inline' or 'attachment'")

	customID := "custom-123"
	acl := "public-read"
	privateACL := "private"
	invalidACL := "invalid-acl"
	contentDisposition := "inline"
	invalidContentDisposition := "invalid-disposition"

	validPayload := &domain.UploadRequest{
		Files: []domain.Files{
			{
				Name:     "test-image.jpg",
				Size:     1024,
				Type:     "image/jpeg",
				CustomID: &customID,
			},
		},
		ACL:                &acl,
		Metadata:           map[string]string{"key": "value"},
		ContentDisposition: &contentDisposition,
	}

	type Given struct {
		payload    *domain.UploadRequest
		mockUpload func(m *client.MockHttpAPI)
	}

	type Expected struct {
		url string
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful upload": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					successResponse := `{
						"data": [
							{
								"data": {
									"key": "abc123",
									"url": "https://utfs.io/f/abc123",
									"appUrl": "https://uploadthing.com/f/abc123",
									"name": "test-image.jpg",
									"size": 1024,
									"customId": "custom-123"
								},
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(successResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "https://utfs.io/f/abc123",
				err: nil,
			},
		},
		"Successful upload with defaults applied": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "test.png",
							Size: 2048,
							Type: "image/png",
						},
					},
				},
				mockUpload: func(m *client.MockHttpAPI) {
					successResponse := `{
						"data": [
							{
								"data": {
									"key": "def456",
									"url": "https://utfs.io/f/def456",
									"appUrl": "https://uploadthing.com/f/def456",
									"name": "test.png",
									"size": 2048
								},
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(successResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "https://utfs.io/f/def456",
				err: nil,
			},
		},
		"Successful upload with private ACL": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "private.jpg",
							Size: 512,
							Type: "image/jpeg",
						},
					},
					ACL: &privateACL,
				},
				mockUpload: func(m *client.MockHttpAPI) {
					successResponse := `{
						"data": [
							{
								"data": {
									"key": "ghi789",
									"url": "https://utfs.io/f/ghi789",
									"appUrl": "https://uploadthing.com/f/ghi789",
									"name": "private.jpg",
									"size": 512
								},
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(successResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "https://utfs.io/f/ghi789",
				err: nil,
			},
		},
		"HTTP client error": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 500,
						}, httpErr)
				},
			},
			expected: Expected{
				url: "",
				err: fmt.Errorf("failed to send HTTP request: %w", httpErr),
			},
		},
		"Non-200 response from server": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 400,
							Status:     "400 Bad Request",
							Body:       io.NopCloser(bytes.NewBufferString(`{"error": "invalid request"}`)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("failed to upload image: status=400 Bad Request message={\"error\": \"invalid request\"}"),
			},
		},
		"Invalid JSON response": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("failed to decode uploadthing response:"),
			},
		},
		"Response with error from UploadThing": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					errorResponse := `{
						"data": [
							{
								"data": null,
								"error": {
									"code": "FILE_TOO_LARGE",
									"message": "File size exceeds limit",
									"data": {}
								}
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(errorResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("invalid uploadthing response: uploadthing: data[0] error: File size exceeds limit (code: FILE_TOO_LARGE)"),
			},
		},
		"Response with missing data": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					invalidResponse := `{
						"data": [
							{
								"data": null,
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(invalidResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("invalid uploadthing response: uploadthing: data[0]: data missing"),
			},
		},
		"Response with missing key": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					invalidResponse := `{
						"data": [
							{
								"data": {
									"key": "",
									"url": "https://utfs.io/f/abc123",
									"appUrl": "https://uploadthing.com/f/abc123",
									"name": "test.jpg",
									"size": 1024
								},
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(invalidResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("invalid uploadthing response: uploadthing: data[0].key missing"),
			},
		},
		"Response with missing URL": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					invalidResponse := `{
						"data": [
							{
								"data": {
									"key": "abc123",
									"url": "",
									"appUrl": "https://uploadthing.com/f/abc123",
									"name": "test.jpg",
									"size": 1024
								},
								"error": null
							}
						]
					}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(invalidResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("invalid uploadthing response: uploadthing: data[0].url missing"),
			},
		},
		"Response with empty data array": {
			given: Given{
				payload: validPayload,
				mockUpload: func(m *client.MockHttpAPI) {
					invalidResponse := `{"data": []}`
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(invalidResponse)),
						}, nil)
				},
			},
			expected: Expected{
				url: "",
				err: errors.New("invalid uploadthing response: uploadthing: response returned no files"),
			},
		},
		"Validation error: missing files": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{},
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: missingFilesErr,
			},
		},
		"Validation error: missing file name": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "",
							Size: 1024,
							Type: "image/jpeg",
						},
					},
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: missingNameErr,
			},
		},
		"Validation error: invalid file size": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "test.jpg",
							Size: 0,
							Type: "image/jpeg",
						},
					},
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: invalidSizeErr,
			},
		},
		"Validation error: missing file type": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "test.jpg",
							Size: 1024,
							Type: "",
						},
					},
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: missingTypeErr,
			},
		},
		"Validation error: invalid ACL": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "test.jpg",
							Size: 1024,
							Type: "image/jpeg",
						},
					},
					ACL: &invalidACL,
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: invalidACLErr,
			},
		},
		"Validation error: invalid content disposition": {
			given: Given{
				payload: &domain.UploadRequest{
					Files: []domain.Files{
						{
							Name: "test.jpg",
							Size: 1024,
							Type: "image/jpeg",
						},
					},
					ContentDisposition: &invalidContentDisposition,
				},
				mockUpload: nil,
			},
			expected: Expected{
				url: "",
				err: invalidContentDispositionErr,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			fixture := newImageRepositoryTestFixture(t)
			if tc.given.mockUpload != nil {
				tc.given.mockUpload(fixture.mockHttpAPI)
			}

			// Act
			url, err := fixture.imageRepository.Upload(context.Background(), tc.given.payload)

			// Assert
			if tc.expected.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.err.Error())
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.url, url)
			}
		})
	}
}

func TestNewImageRepository(t *testing.T) {
	t.Run("Creates repository with provided httpAPI", func(t *testing.T) {
		// Arrange
		mockHttpAPI := client.NewMockHttpAPI(t)
		cfg := ImageRepositoryConfig{
			UploadthingToken: "test_token",
			httpAPI:          mockHttpAPI,
		}

		// Act
		repo := NewImageRepository(cfg)

		// Assert
		assert.NotNil(t, repo)
		concreteRepo, ok := repo.(*imageRepository)
		assert.True(t, ok)
		assert.Equal(t, "test_token", concreteRepo.uploadthingToken)
		assert.Equal(t, mockHttpAPI, concreteRepo.httpAPI)
	})

	t.Run("Creates repository with default httpAPI when nil", func(t *testing.T) {
		// Arrange
		cfg := ImageRepositoryConfig{
			UploadthingToken: "test_token",
			httpAPI:          nil,
		}

		// Act
		repo := NewImageRepository(cfg)

		// Assert
		assert.NotNil(t, repo)
		concreteRepo, ok := repo.(*imageRepository)
		assert.True(t, ok)
		assert.Equal(t, "test_token", concreteRepo.uploadthingToken)
		assert.NotNil(t, concreteRepo.httpAPI)
	})
}
