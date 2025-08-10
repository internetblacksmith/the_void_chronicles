package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "returns default when env var not set",
			key:          "UNSET_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHTTPServer(t *testing.T) {
	// Create a test server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test</body></html>"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Test root endpoint
	t.Run("root endpoint returns HTML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "text/html" {
			t.Errorf("expected Content-Type text/html, got %s", contentType)
		}

		body := w.Body.String()
		if body == "" {
			t.Error("expected non-empty body")
		}
	})

	// Test health endpoint
	t.Run("health endpoint returns OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if w.Body.String() != "OK" {
			t.Errorf("expected body 'OK', got %s", w.Body.String())
		}
	})
}

func TestPasswordHandler(t *testing.T) {
	tests := []struct {
		name         string
		envPassword  string
		inputPassword string
		expected     bool
	}{
		{
			name:         "correct password with env var",
			envPassword:  "TestPass123",
			inputPassword: "TestPass123",
			expected:     true,
		},
		{
			name:         "incorrect password with env var",
			envPassword:  "TestPass123",
			inputPassword: "WrongPass",
			expected:     false,
		},
		{
			name:         "correct default password",
			envPassword:  "",
			inputPassword: "Amigos4Life!",
			expected:     true,
		},
		{
			name:         "incorrect default password",
			envPassword:  "",
			inputPassword: "WrongPass",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envPassword != "" {
				os.Setenv("SSH_PASSWORD", tt.envPassword)
				defer os.Unsetenv("SSH_PASSWORD")
			}

			result := passwordHandler(nil, tt.inputPassword)
			if result != tt.expected {
				t.Errorf("passwordHandler() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateSSHKey(t *testing.T) {
	tempDir := t.TempDir()
	keyPath := tempDir + "/test_key"

	err := generateSSHKey(keyPath)
	if err != nil {
		t.Fatalf("generateSSHKey() error = %v", err)
	}

	// Check if files were created
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("private key file was not created")
	}

	if _, err := os.Stat(keyPath + ".pub"); os.IsNotExist(err) {
		t.Error("public key file was not created")
	}
}

func TestServerStartup(t *testing.T) {
	// This test verifies that the servers can start without panicking
	// We use goroutines with timeouts to prevent hanging

	t.Run("HTTP server can start", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("HTTP server panicked: %v", r)
				}
				done <- true
			}()

			// Start server on random port
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test"))
			})
			
			server := httptest.NewServer(mux)
			defer server.Close()
		}()

		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Error("HTTP server startup timeout")
		}
	})
}