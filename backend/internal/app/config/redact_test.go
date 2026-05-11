package config

import (
	"testing"
)

func TestRedisURL_Redacted(t *testing.T) {
	const expectedRedacted = "%2A%2A%2AREDACTED%2A%2A%2A"

	tests := []struct {
		name string
		url  RedisURL
		want RedisURL
	}{
		{
			name: "URL with password should be redacted",
			url:  "redis://user:password123@localhost:6379",
			want: "redis://user:" + expectedRedacted + "@localhost:6379",
		},
		{
			name: "URL with only username should be redacted",
			url:  "redis://admin@localhost:6379",
			want: "redis://admin:" + expectedRedacted + "@localhost:6379",
		},
		{
			name: "URL without user info should remain unchanged",
			url:  "redis://localhost:6379",
			want: "redis://localhost:6379",
		},
		{
			name: "Invalid URL should be returned as is",
			url:  "://invalid-url",
			want: "://invalid-url",
		},
		{
			name: "Empty URL should return empty string",
			url:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.Redacted(); got != tt.want {
				t.Errorf("RedisURL.Redacted() = %v, want %v", got, tt.want)
			}
		})
	}
}
