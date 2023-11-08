package oidconfig_test

import (
	"testing"

	"github.com/numberly/gangway/internal/config"
	"github.com/numberly/gangway/internal/oidconfig"
	"github.com/stretchr/testify/assert"
)

func TestInitOIDCProviders(t *testing.T) {
	// Mock configuration
	clusterCfg := &config.MultiClusterConfig{
		Host: "localhost",
		Port: 8080,
		Clusters: map[string][]config.Config{
			"Development": {
				{
					EnvPrefix:     "dev",
					ClusterName:   "example-cluster",
					ProviderURL:   "https://accounts.google.com",
					ClientID:      "client-id",
					ClientSecret:  "client-secret",
					RedirectURL:   "https://example.com/callback",
					Scopes:        []string{"openid", "profile", "email"},
					UsernameClaim: "preferred_username",
					EmailClaim:    "email",
					APIServerURL:  "https://api.example.com",
				},
			},
		},
		HTTPPath: "/",
		// ... other fields
	}

	// Initialize OIDC providers
	err := oidconfig.InitOIDCProviders(*clusterCfg)
	assert.NoError(t, err)

	// Check if the provider configuration is stored
	config, ok := oidconfig.GetOIDCProviderConfig("example-cluster")
	assert.True(t, ok)
	assert.NotNil(t, config.Provider)
	assert.NotNil(t, config.OAuth2Config)
	assert.NotNil(t, config.Verifier)
}

func TestGetOIDCProviderConfig_NotFound(t *testing.T) {
	// Attempt to retrieve a non-existent provider configuration
	_, ok := oidconfig.GetOIDCProviderConfig("non-existent-cluster")
	assert.False(t, ok)
}
