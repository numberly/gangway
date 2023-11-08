package oidconfig

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/numberly/gangway/internal/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Structure pour stocker le provider OIDC et la configuration OAuth2.
type OIDCProviderConfig struct {
	Provider     *oidc.Provider
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

// Map pour stocker les configurations par clusterName.
var ClusterOIDCConfigs = make(map[string]*OIDCProviderConfig)

func InitOIDCProviders(multiClusterConfig config.MultiClusterConfig) error {
	ctx := context.Background()

	for _, clusterConfigs := range multiClusterConfig.Clusters {
		for _, clusterConfig := range clusterConfigs {
			// Créer le provider OIDC.
			provider, err := oidc.NewProvider(ctx, clusterConfig.ProviderURL)
			if err != nil {
				log.Errorf("Could not create OIDC provider for cluster %s: %v", clusterConfig.ClusterName, err)
				return err
			}

			// Créer le verifier OIDC.
			verifier := provider.Verifier(&oidc.Config{ClientID: clusterConfig.ClientID})

			// Configurer le client OAuth2.
			oauth2Cfg := &oauth2.Config{
				ClientID:     clusterConfig.ClientID,
				ClientSecret: clusterConfig.ClientSecret,
				RedirectURL:  clusterConfig.RedirectURL,
				Scopes:       clusterConfig.Scopes,
				Endpoint:     provider.Endpoint(),
			}

			// Stocker la configuration dans la map.
			ClusterOIDCConfigs[clusterConfig.ClusterName] = &OIDCProviderConfig{
				Provider:     provider,
				OAuth2Config: oauth2Cfg,
				Verifier:     verifier,
			}
		}
	}

	return nil
}

func GetOIDCProviderConfig(clusterName string) (*OIDCProviderConfig, bool) {
	log.Printf("Voici le contenu de ClusterOIDCConfigs : %v", ClusterOIDCConfigs[clusterName])
	config, ok := ClusterOIDCConfigs[clusterName]
	return config, ok
}
