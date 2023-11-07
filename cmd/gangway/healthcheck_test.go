package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/numberly/gangway/internal/config"
)

func TestLivenessHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/live", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(livenessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestReadinessHandler(t *testing.T) {
	clusterCfg = &config.MultiClusterConfig{
		Host: "localhost",
		Port: 8080,
		Clusters: map[string][]config.Config{
			"Development": {
				{
					EnvPrefix:              "dev",
					ClusterName:            "example-cluster",
					ProviderURL:            "https://provider.example.com",
					ClientID:               "client-id",
					ClientSecret:           "client-secret",
					AllowEmptyClientSecret: false,
					Audience:               "https://example.com",
					RedirectURL:            "https://example.com/callback",
					Scopes:                 []string{"openid", "profile", "email"},
					UsernameClaim:          "preferred_username",
					EmailClaim:             "email",
					APIServerURL:           "https://api.example.com",
					ShowClaims:             false,
				},
			},
		},
		HTTPPath: "/",
		// ... other fields
	}
	// Setup a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readinessHandler)

	// Call the handler, which in turn calls clustersHomeCheck.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ready": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
