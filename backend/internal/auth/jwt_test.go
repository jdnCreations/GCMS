package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	uuid, err := uuid.Parse("bd7d92e2-a631-4258-8b58-db200b311153")
	if err != nil {
		t.Fatalf("Could not parse uuid")
	}
	jwt, err := MakeJWT(uuid, "test", 1 * time.Minute)
	if err != nil {
		t.Fatalf("Could not create JWT, err: %s", err.Error())
	}

	jwt2, err := MakeJWT(uuid, "test", -1 * time.Second)
	if err != nil {
		t.Fatalf("Could not create JWT, err: %s", err.Error())
	}

	validated, err := ValidateJWT(jwt, "test")
	if err != nil {
		t.Errorf("Could not validate JWT, err: %s", err.Error())
	}

	_, err = ValidateJWT(jwt2, "test")
	if err == nil {
		t.Errorf("Expected error for expired token, got none")
	}

	if validated.String() != uuid.String() {
		t.Errorf("Expected UUID %s, got %s", uuid, validated)
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name string
		headerValue string
		expectedToken string
		expectError bool
	}{
		{
			name: "valid bearer token",
			headerValue: "Bearer abc123",
			expectedToken: "abc123",
			expectError: false,
		},
		{
			name: "missing bearer token",
			headerValue: "",
			expectedToken: "",
			expectError: true,
		},
		{
			name: "extra spaces",
			headerValue: "Bearer    abc123",
			expectedToken: "abc123",
			expectError: false,
		},
		{
			name: "lowercase bearer",
			headerValue: "bearer abc123",
			expectedToken: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerValue != "" {
				headers.Add("Authorization", tt.headerValue)
			}
			token, err := GetBearerToken(headers)
			
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}