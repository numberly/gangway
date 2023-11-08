package jwt_test

import (
	"testing"
	"time"

	"github.com/numberly/gangway/internal/config"
	"github.com/numberly/gangway/internal/jwt"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	clusterName := "test-cluster"
	oauth2State := "state-123"
	clusterCfg := config.MultiClusterConfig{
		SessionSecurityKey: "supersecretkey",
	}

	tokenString, err := jwt.CreateToken(clusterName, oauth2State, clusterCfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Validate the token
	claims, err := jwt.ValidateToken(tokenString, clusterCfg)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, clusterName, claims.ClusterName)
	assert.Equal(t, oauth2State, claims.OAuth2State)
	assert.WithinDuration(t, time.Unix(claims.ExpiresAt, 0), time.Now().Add(1*time.Hour), 5*time.Second)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	clusterCfg := config.MultiClusterConfig{
		SessionSecurityKey: "supersecretkey",
	}
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbHVzdGVyX25hbWUiOiJ0ZXN0LWNsdXN0ZXIiLCJvYXV0aDJfc3RhdGUiOiJzdGF0ZS0xMjMifQ.WrongSignature"

	_, err := jwt.ValidateToken(invalidToken, clusterCfg)
	assert.Error(t, err)
}

func TestUpdateToken(t *testing.T) {
	clusterName := "test-cluster"
	oauth2State := "state-123"
	clusterCfg := config.MultiClusterConfig{
		SessionSecurityKey: "supersecretkey",
	}
	tokenId := "token-123"
	refreshTokenId := "refresh-123"

	originalTokenString, _ := jwt.CreateToken(clusterName, oauth2State, clusterCfg)
	updatedTokenString, err := jwt.UpdateToken(originalTokenString, tokenId, refreshTokenId, clusterCfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedTokenString)

	// Validate the updated token
	updatedClaims, err := jwt.ValidateToken(updatedTokenString, clusterCfg)
	assert.NoError(t, err)
	assert.Equal(t, tokenId, updatedClaims.OAuth2TokenId)
	assert.Equal(t, refreshTokenId, updatedClaims.OAuth2RefreshId)
}
