package config

import (
	"os"
	"testing"
)

// createTempFileWithContent crée un fichier temporaire avec le contenu spécifié et retourne son chemin.
func createTempFileWithContent(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

// removeTempFile supprime le fichier temporaire spécifié.
func removeTempFile(t *testing.T, filename string) {
	t.Helper()
	if err := os.Remove(filename); err != nil {
		t.Fatal(err)
	}
}

// TestNewMultiClusterConfig tests the NewMultiClusterConfig function for various scenarios.
func TestNewMultiClusterConfig(t *testing.T) {
	t.Run("Test with valid config file", func(t *testing.T) {
		content := `
httpPath: "/test"
serveTLS: true
clusters:
  prod:
    - clusterName: "prod-cluster"
      providerURL: "https://provider.example.com"
      clientID: "my-client-id"
      clientSecret: "my-client-secret"
      redirectURL: "https://redirect.example.com"
      apiServerURL: "https://apiserver.example.com"
`
		filename := createTempFileWithContent(t, content)
		defer removeTempFile(t, filename)

		cfg, err := NewMultiClusterConfig(filename)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.HTTPPath != "/test" {
			t.Errorf("expected HTTPPath to be '/test', got %s", cfg.HTTPPath)
		}
		if !cfg.ServeTLS {
			t.Errorf("expected ServeTLS to be true, got %v", cfg.ServeTLS)
		}
		if len(cfg.Clusters["prod"]) == 0 {
			t.Errorf("expected at least one cluster in 'prod', got %v", cfg.Clusters["prod"])
		}
	})

	t.Run("Test with missing mandatory fields", func(t *testing.T) {
		content := `
clusters:
  prod:
    - clientID: "my-client-id"
`
		filename := createTempFileWithContent(t, content)
		defer removeTempFile(t, filename)

		_, err := NewMultiClusterConfig(filename)
		if err == nil {
			t.Fatalf("expected an error due to missing mandatory fields, got none")
		}
	})

	t.Run("Test with environment variables", func(t *testing.T) {
		content := `clusters: {}` // Empty clusters to load from envvars
		filename := createTempFileWithContent(t, content)
		defer removeTempFile(t, filename)

		// Set environment variables for testing
		os.Setenv("GANGWAY_CLUSTER_NAME", "env-cluster")
		os.Setenv("GANGWAY_PROVIDER_URL", "https://envprovider.example.com")
		os.Setenv("GANGWAY_CLIENT_ID", "env-client-id")
		os.Setenv("GANGWAY_CLIENT_SECRET", "env-client-secret")
		os.Setenv("GANGWAY_REDIRECT_URL", "https://envredirect.example.com")
		os.Setenv("GANGWAY_APISERVER_URL", "https://envapiserver.example.com")
		defer func() {
			// Clean up environment variables
			os.Unsetenv("GANGWAY_CLUSTER_NAME")
			os.Unsetenv("GANGWAY_PROVIDER_URL")
			os.Unsetenv("GANGWAY_CLIENT_ID")
			os.Unsetenv("GANGWAY_CLIENT_SECRET")
			os.Unsetenv("GANGWAY_REDIRECT_URL")
			os.Unsetenv("GANGWAY_APISERVER_URL")
		}()

		cfg, err := NewMultiClusterConfig(filename)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.Clusters == nil {
			t.Fatalf("expected clusters to be loaded from environment, got nil")
		}
	})
}
