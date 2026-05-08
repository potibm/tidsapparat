package config

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisURL_URLObject(t *testing.T) {
	tests := []struct {
		name     string
		input    *RedisURL
		expected string
	}{
		{
			name:     "valid redis URL",
			input:    redisURLPtr("redis://localhost:6379/0"),
			expected: "redis://localhost:6379/0",
		},
		{
			name:     "nil receiver",
			input:    nil,
			expected: "",
		},
		{
			name:     "invalid URL",
			input:    redisURLPtr("://invalid"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.URLObject()
			if tt.expected == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestRedisURL_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    *RedisURL
		expected bool
	}{
		{
			name:     "valid redis URL",
			input:    redisURLPtr("redis://localhost:6379"),
			expected: true,
		},
		{
			name:     "nil receiver",
			input:    nil,
			expected: true,
		},
		{
			name:     "invalid URL",
			input:    redisURLPtr("://invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestRedisURL_RedisOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    RedisURL
		expected *redis.Options
	}{
		{
			name:  "full URL with password and DB",
			input: "redis://:secret@localhost:6379/3",
			expected: &redis.Options{
				Addr:     "localhost:6379",
				Password: "secret",
				DB:       3,
			},
		},
		{
			name:  "URL without password",
			input: "redis://localhost:6379/1",
			expected: &redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       1,
			},
		},
		{
			name:  "URL without DB defaults to 0",
			input: "redis://localhost:6379",
			expected: &redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		{
			name:  "empty path defaults to DB 0",
			input: "redis://localhost:6379/",
			expected: &redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		{
			name:  "invalid DB path falls back to 0",
			input: "redis://localhost:6379/abc",
			expected: &redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		{
			name:     "nil URL object returns nil",
			input:    "://invalid",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.RedisOptions()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func redisURLPtr(s string) *RedisURL {
	r := RedisURL(s)

	return &r
}
