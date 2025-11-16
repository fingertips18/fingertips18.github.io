package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type imageHandlerTestFixture struct {
	t             *testing.T
	mockImageRepo *mockRepo.MockImageRepository
	imageHandler  ImageHandler
}

func newImageHandlerTestFixture(t *testing.T) *imageHandlerTestFixture {
	mockImageRepo := new(mockRepo.MockImageRepository)

	imageHandler := NewImageServiceHandler(
		ImageServiceConfig{
			imageRepo: mockImageRepo,
		},
	)

	return &imageHandlerTestFixture{
		t:             t,
		mockImageRepo: mockImageRepo,
		imageHandler:  imageHandler,
	}
}

func TestImageServiceHandler_Upload(t *testing.T) {
	validFile := FilesDTO{
		Name: "profile.jpg",
		Size: 1024,
		Type: "image/jpeg",
	}

	validReq := UploadRequestDTO{
		Files: []FilesDTO{validFile},
	}
	validBody, _ := json.Marshal(validReq)

	expectedURL := "https://uploadthing.com/f/abc123.jpg"
	validResp, _ := json.Marshal(UploadResponseDTO{URL: expectedURL})

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockImageRepository)
	}
	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 &&
								req.Files[0].Name == "profile.jpg" &&
								req.Files[0].Size == 1024 &&
								req.Files[0].Type == "image/jpeg"
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"invalid method": {
			given: Given{
				method: http.MethodGet,
				body:   "",
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only POST is supported\n",
			},
		},
		"invalid json": {
			given: Given{
				method: http.MethodPost,
				body:   `{"files":}`,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"empty body": {
			given: Given{
				method: http.MethodPost,
				body:   "",
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.AnythingOfType("*domain.UploadRequest")).
						Return("", errors.New("upload failed"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to upload image: upload failed\n",
			},
		},
		"multiple files": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{Name: "image1.jpg", Size: 1024, Type: "image/jpeg"},
							{Name: "image2.png", Size: 2048, Type: "image/png"},
							{Name: "image3.gif", Size: 512, Type: "image/gif"},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 3 &&
								req.Files[0].Name == "image1.jpg" &&
								req.Files[1].Name == "image2.png" &&
								req.Files[2].Name == "image3.gif"
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"empty files array": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 0
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"with all optional fields": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					acl := "public-read"
					contentDisposition := "inline"
					req := UploadRequestDTO{
						Files:              []FilesDTO{validFile},
						ACL:                &acl,
						ContentDisposition: &contentDisposition,
						Metadata: map[string]string{
							"user_id":    "123",
							"project_id": "456",
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return req.ACL != nil &&
								*req.ACL == "public-read" &&
								req.ContentDisposition != nil &&
								*req.ContentDisposition == "inline" &&
								req.Metadata != nil
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"with custom_id": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					customID := "custom-file-id-123"
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{
								Name:     "profile.jpg",
								Size:     1024,
								Type:     "image/jpeg",
								CustomID: &customID,
							},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 &&
								req.Files[0].CustomID != nil &&
								*req.Files[0].CustomID == "custom-file-id-123"
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"large file size": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{
								Name: "large-image.jpg",
								Size: 104857600, // 100MB
								Type: "image/jpeg",
							},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 && req.Files[0].Size == 104857600
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"unicode filename": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{
								Name: "画像ファイル-🖼️.jpg",
								Size: 1024,
								Type: "image/jpeg",
							},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 &&
								strings.Contains(req.Files[0].Name, "画像ファイル")
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"very long filename": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{
								Name: strings.Repeat("a", 500) + ".jpg",
								Size: 1024,
								Type: "image/jpeg",
							},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 && len(req.Files[0].Name) > 500
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"various image types": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{Name: "image.jpg", Size: 1024, Type: "image/jpeg"},
							{Name: "image.png", Size: 2048, Type: "image/png"},
							{Name: "image.gif", Size: 512, Type: "image/gif"},
							{Name: "image.webp", Size: 256, Type: "image/webp"},
							{Name: "image.svg", Size: 128, Type: "image/svg+xml"},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 5
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"malformed JSON with extra fields": {
			given: Given{
				method: http.MethodPost,
				body:   `{"files":[{"name":"test.jpg","size":1024,"type":"image/jpeg"}],"extra_field":"ignored"}`,
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.AnythingOfType("*domain.UploadRequest")).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"zero size file": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{
							{
								Name: "empty.jpg",
								Size: 0,
								Type: "image/jpeg",
							},
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return len(req.Files) == 1 && req.Files[0].Size == 0
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
		"complex metadata": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					req := UploadRequestDTO{
						Files: []FilesDTO{validFile},
						Metadata: map[string]interface{}{
							"user": map[string]string{
								"id":   "123",
								"name": "John Doe",
							},
							"tags":      []string{"profile", "avatar", "user"},
							"timestamp": 1234567890,
							"public":    true,
						},
					}
					b, _ := json.Marshal(req)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockImageRepository) {
					m.EXPECT().
						Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
							return req.Metadata != nil
						})).
						Return(expectedURL, nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: string(validResp),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newImageHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockImageRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/image/upload", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.imageHandler.Upload(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockImageRepo.AssertExpectations(t)
		})
	}
}

func TestImageServiceHandler_Upload_Routing(t *testing.T) {
	validFile := FilesDTO{
		Name: "profile.jpg",
		Size: 1024,
		Type: "image/jpeg",
	}

	validReq := UploadRequestDTO{
		Files: []FilesDTO{validFile},
	}
	validBody, _ := json.Marshal(validReq)

	expectedURL := "https://uploadthing.com/f/abc123.jpg"
	expectedResp, _ := json.Marshal(UploadResponseDTO{URL: expectedURL})

	f := newImageHandlerTestFixture(t)

	// Mock expectation
	f.mockImageRepo.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
			return len(req.Files) == 1 &&
				req.Files[0].Name == "profile.jpg"
		})).
		Return(expectedURL, nil)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/image/upload", bytes.NewReader(validBody))
	w := httptest.NewRecorder()

	// Verify handler implements http.Handler
	handler, ok := f.imageHandler.(http.Handler)
	assert.True(t, ok, "imageHandler should implement http.Handler")

	// Route through ServeHTTP
	handler.ServeHTTP(w, req)

	// Validate response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockImageRepo.AssertExpectations(t)
}

func TestImageServiceHandler_ServeHTTP_NotFound(t *testing.T) {
	f := newImageHandlerTestFixture(t)

	tests := []struct {
		name string
		path string
	}{
		{"root path", "/image"},
		{"invalid path", "/image/invalid"},
		{"nested path", "/image/upload/extra"},
		{"different path", "/image/download"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler, ok := f.imageHandler.(http.Handler)
			assert.True(t, ok, "imageHandler should implement http.Handler")

			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		})
	}
}

func TestImageServiceHandler_ServeHTTP_TrailingSlash(t *testing.T) {
	validFile := FilesDTO{
		Name: "profile.jpg",
		Size: 1024,
		Type: "image/jpeg",
	}

	validReq := UploadRequestDTO{
		Files: []FilesDTO{validFile},
	}
	validBody, _ := json.Marshal(validReq)

	expectedURL := "https://uploadthing.com/f/abc123.jpg"
	expectedResp, _ := json.Marshal(UploadResponseDTO{URL: expectedURL})

	f := newImageHandlerTestFixture(t)

	// Mock expectation
	f.mockImageRepo.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(req *domain.UploadRequest) bool {
			return len(req.Files) == 1
		})).
		Return(expectedURL, nil)

	// Create request with trailing slash
	req := httptest.NewRequest(http.MethodPost, "/image/upload/", bytes.NewReader(validBody))
	w := httptest.NewRecorder()

	handler, ok := f.imageHandler.(http.Handler)
	assert.True(t, ok, "imageHandler should implement http.Handler")

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockImageRepo.AssertExpectations(t)
}
