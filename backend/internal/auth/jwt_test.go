package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	uuid, err := uuid.Parse("9d961878-8641-4ba9-b194-9b9b48122384")
	if err != nil {
		t.Fatalf("Could not parse uuid")
	}
	jwt, err := MakeJWT(uuid, "test")
	if err != nil {
		t.Fatalf("Could not create JWT, err: %s", err.Error())
	}

	validated, err := ValidateJWT(jwt, "test")
	if err != nil {
		t.Errorf("Could not validate JWT, err: %s", err.Error())
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