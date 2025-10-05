package utils

import (
	"net/url"
	"testing"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetQueryInt32(t *testing.T) {
	tests := map[string]struct {
		q    url.Values
		key  string
		def  int32
		want int32
	}{
		"missing key returns default": {
			q:    url.Values{},
			key:  "page",
			def:  5,
			want: 5,
		},
		"valid integer string": {
			q:    url.Values{"page": {"10"}},
			key:  "page",
			def:  1,
			want: 10,
		},
		"invalid integer string": {
			q:    url.Values{"page": {"abc"}},
			key:  "page",
			def:  2,
			want: 2,
		},
		"large integer within int32 range": {
			q:    url.Values{"page": {"2147483647"}},
			key:  "page",
			def:  0,
			want: 2147483647,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetQueryInt32(tt.q, tt.key, tt.def)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetQuerySortBy(t *testing.T) {
	tests := map[string]struct {
		q      url.Values
		key    string
		want   *domain.SortBy
		hasErr bool
	}{
		"missing key returns nil": {
			q:      url.Values{},
			key:    "sort_by",
			want:   nil,
			hasErr: false,
		},
		"valid CreatedAt": {
			q:      url.Values{"sort_by": {string(domain.CreatedAt)}},
			key:    "sort_by",
			want:   func() *domain.SortBy { s := domain.CreatedAt; return &s }(),
			hasErr: false,
		},
		"valid UpdatedAt": {
			q:      url.Values{"sort_by": {string(domain.UpdatedAt)}},
			key:    "sort_by",
			want:   func() *domain.SortBy { s := domain.UpdatedAt; return &s }(),
			hasErr: false,
		},
		"invalid value": {
			q:      url.Values{"sort_by": {"invalid"}},
			key:    "sort_by",
			want:   nil,
			hasErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetQuerySortBy(tt.q, tt.key)
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetQueryBool(t *testing.T) {
	tests := map[string]struct {
		q    url.Values
		key  string
		def  bool
		want bool
	}{
		"missing key returns default": {
			q:    url.Values{},
			key:  "active",
			def:  true,
			want: true,
		},
		"valid true": {
			q:    url.Values{"active": {"true"}},
			key:  "active",
			def:  false,
			want: true,
		},
		"valid false": {
			q:    url.Values{"active": {"false"}},
			key:  "active",
			def:  true,
			want: false,
		},
		"truthy 1": {
			q:    url.Values{"active": {"1"}},
			key:  "active",
			def:  false,
			want: true,
		},
		"invalid bool string returns default": {
			q:    url.Values{"active": {"notabool"}},
			key:  "active",
			def:  true,
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetQueryBool(tt.q, tt.key, tt.def)
			assert.Equal(t, tt.want, got)
		})
	}
}
