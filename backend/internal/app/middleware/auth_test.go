package middleware

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildHTTPClient(t *testing.T) {
	t.Run("with TLS verification", func(t *testing.T) {
		client, err := buildHTTPClient(false)
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Equal(t, 10*time.Second, client.Timeout)

		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)

		if transport.TLSClientConfig != nil {
			assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
		}
	})

	t.Run("skip TLS verification", func(t *testing.T) {
		client, err := buildHTTPClient(true)
		require.NoError(t, err)
		require.NotNil(t, client)

		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)
		require.NotNil(t, transport.TLSClientConfig)
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing header",
			authHeader:  "",
			wantErr:     true,
			errContains: "missing Authorization header",
		},
		{
			name:       "valid Bearer token",
			authHeader: "Bearer my-token-123",
			wantToken:  "my-token-123",
		},
		{
			name:       "lowercase bearer",
			authHeader: "bearer my-token-123",
			wantToken:  "my-token-123",
		},
		{
			name:        "wrong scheme",
			authHeader:  "Basic dXNlcjpwYXNz",
			wantErr:     true,
			errContains: "invalid token format",
		},
		{
			name:        "only Bearer prefix",
			authHeader:  "Bearer",
			wantErr:     true,
			errContains: "invalid token format",
		},
		{
			name:        "three parts",
			authHeader:  "Bearer token extra",
			wantErr:     true,
			errContains: "invalid token format",
		},
		{
			name:       "mixed case prefix",
			authHeader: "BeArEr my-token",
			wantToken:  "my-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := extractBearerToken(tt.authHeader)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func TestValidateTokenAndGetUserID(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	givenKey := keyfunc.NewGivenRSA(&privateKey.PublicKey, keyfunc.GivenKeyOptions{
		Algorithm: "RS256",
	})
	jwks := keyfunc.NewGiven(map[string]keyfunc.GivenKey{"test-key": givenKey})

	expectedIssuer := "https://issuer.example.com"
	expectedUserID := "user-42"

	makeToken := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "test-key"
		s, err := token.SignedString(privateKey)
		require.NoError(t, err)

		return s
	}

	t.Run("valid token", func(t *testing.T) {
		tokenString := makeToken(jwt.MapClaims{
			"iss": expectedIssuer,
			"sub": expectedUserID,
		})

		userID, err := validateTokenAndGetUserID(tokenString, jwks, expectedIssuer)
		require.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("invalid token string", func(t *testing.T) {
		_, err := validateTokenAndGetUserID("not-a-token", jwks, expectedIssuer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("wrong issuer", func(t *testing.T) {
		tokenString := makeToken(jwt.MapClaims{
			"iss": "https://wrong-issuer.example.com",
			"sub": expectedUserID,
		})

		_, err := validateTokenAndGetUserID(tokenString, jwks, expectedIssuer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid issuer")
	})

	t.Run("missing issuer claim", func(t *testing.T) {
		tokenString := makeToken(jwt.MapClaims{
			"sub": expectedUserID,
		})

		_, err := validateTokenAndGetUserID(tokenString, jwks, expectedIssuer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid issuer")
	})

	t.Run("missing subject", func(t *testing.T) {
		tokenString := makeToken(jwt.MapClaims{
			"iss": expectedIssuer,
		})

		_, err := validateTokenAndGetUserID(tokenString, jwks, expectedIssuer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing subject")
	})

	t.Run("empty subject", func(t *testing.T) {
		tokenString := makeToken(jwt.MapClaims{
			"iss": expectedIssuer,
			"sub": "",
		})

		_, err := validateTokenAndGetUserID(tokenString, jwks, expectedIssuer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing subject")
	})
}

func TestInitJWKSSuccess(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwkData := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(privateKey.N.Bytes()),
				"e":   "AQAB",
			},
		},
	}
	jwksJSON, err := json.Marshal(jwkData)
	require.NoError(t, err)

	issuer := "https://test-issuer.example.com"

	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(jwksJSON)
		require.NoError(t, err)
	}))
	defer jwksServer.Close()

	discovery := OIDCDiscovery{
		JwksURI: jwksServer.URL,
		Issuer:  issuer,
	}
	discoveryJSON, err := json.Marshal(discovery)
	require.NoError(t, err)

	discoveryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/openid-configuration" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(discoveryJSON)
		require.NoError(t, err)
	}))
	defer discoveryServer.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	jwks, expectedIssuer, err := initJWKS(context.Background(), client, discoveryServer.URL)
	require.NoError(t, err)
	require.NotNil(t, jwks)
	assert.Equal(t, issuer, expectedIssuer)
}

func TestInitJWKSFailures(t *testing.T) {
	t.Run("non-OK discovery status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := &http.Client{Timeout: 5 * time.Second}
		_, _, err := initJWKS(context.Background(), client, server.URL)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "OIDC discovery returned status")
	})

	t.Run("invalid discovery JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{invalid json`))
			require.NoError(t, err)
		}))
		defer server.Close()

		client := &http.Client{Timeout: 5 * time.Second}
		_, _, err := initJWKS(context.Background(), client, server.URL)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error decoding OIDC discovery document")
	})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwkData := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(privateKey.N.Bytes()),
				"e":   "AQAB",
			},
		},
	}
	jwksJSON, err := json.Marshal(jwkData)
	require.NoError(t, err)

	issuer := "https://test-issuer.example.com"
	userID := "user-123"

	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(jwksJSON)
		require.NoError(t, err)
	}))
	defer jwksServer.Close()

	discovery := OIDCDiscovery{
		JwksURI: jwksServer.URL,
		Issuer:  issuer,
	}
	discoveryJSON, err := json.Marshal(discovery)
	require.NoError(t, err)

	discoveryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(discoveryJSON)
			require.NoError(t, err)

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer discoveryServer.Close()

	makeToken := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "test-key"
		s, err := token.SignedString(privateKey)
		require.NoError(t, err)

		return s
	}

	middleware, err := AuthMiddleware(context.Background(), discoveryServer.URL, "test-client", false)
	require.NoError(t, err)

	t.Run("valid token sets userID in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer "+makeToken(jwt.MapClaims{
			"iss": issuer,
			"sub": userID,
		}))
		c.Request = req

		var capturedUserID string

		var capturedCtxUserID interface{}

		c.Next()

		middleware(c)

		if !c.IsAborted() {
			val, _ := c.Get("userID")
			capturedUserID, _ = val.(string)
			capturedCtxUserID = c.Request.Context().Value(domain.UserIDKey)
		}

		// Re-run with a handler that runs after middleware
		w = httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		req = httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer "+makeToken(jwt.MapClaims{
			"iss": issuer,
			"sub": userID,
		}))
		c.Request = req

		engine.Use(middleware)
		engine.GET("/protected", func(c *gin.Context) {
			val, _ := c.Get("userID")
			capturedUserID, _ = val.(string)
			capturedCtxUserID = c.Request.Context().Value(domain.UserIDKey)
			c.String(http.StatusOK, "ok")
		})
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, userID, capturedUserID)
		assert.Equal(t, userID, capturedCtxUserID)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		c.Request = req

		engine.Use(middleware)
		engine.GET("/protected", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
		c.Request = req

		engine.Use(middleware)
		engine.GET("/protected", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid or expired token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, engine := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
		req.Header.Set("Authorization", "Bearer invalid-token-string")
		c.Request = req

		engine.Use(middleware)
		engine.GET("/protected", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
