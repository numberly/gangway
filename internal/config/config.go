// Copyright © 2017 Heptio
// Copyright © 2017 Craig Tracey
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"sigs.k8s.io/yaml"
)

const hardCodedDefaultSalt = "MkmfuPNHnZBBivy0L0aW"

// Config the configuration field for gangway
type Config struct {
	EnvPrefix              string   `yaml:"-"`
	ClusterName            string   `yaml:"clusterName" envconfig:"cluster_name"`
	ProviderURL            string   `yaml:"providerURL" envconfig:"provider_url"`
	ClientID               string   `yaml:"clientID" envconfig:"client_id"`
	ClientSecret           string   `yaml:"clientSecret" envconfig:"client_secret"`
	AllowEmptyClientSecret bool     `yaml:"allowEmptyClientSecret" envconfig:"allow_empty_client_secret"`
	Audience               string   `yaml:"audience" envconfig:"audience"`
	RedirectURL            string   `yaml:"redirectURL" envconfig:"redirect_url"`
	Scopes                 []string `yaml:"scopes" envconfig:"scopes"`
	UsernameClaim          string   `yaml:"usernameClaim" envconfig:"username_claim"`
	EmailClaim             string   `yaml:"emailClaim" envconfig:"email_claim"`
	APIServerURL           string   `yaml:"apiServerURL" envconfig:"apiserver_url"`
	ClusterCAPath          string   `yaml:"clusterCAPath" envconfig:"cluster_ca_path"`
	ClusterCA              []byte
	ShowClaims             bool `yaml:"showClaims" envconfig:"show_claims"`
}

type MultiClusterConfig struct {
	Host                   string              `yaml:"host"`
	Port                   int                 `yaml:"port"`
	Clusters               map[string][]Config `yaml:"clusters"`
	HTTPPath               string              `yaml:"httpPath" envconfig:"http_path"`
	SessionSecurityKey     string              `yaml:"sessionSecurityKey" envconfig:"session_security_key"`
	SessionSalt            string              `yaml:"sessionSalt" envconfig:"session_salt"`
	CustomHTMLTemplatesDir string              `yaml:"customHTMLTemplatesDir" envconfig:"custom_html_templates_dir"`
	CustomAssetsDir        string              `yaml:"customAssetsDir" envconfig:"custom_assets_dir"`
	ServeTLS               bool                `yaml:"serveTLS" envconfig:"serve_tls"`
	CertFile               string              `yaml:"certFile" envconfig:"cert_file"`
	KeyFile                string              `yaml:"keyFile" envconfig:"key_file"`
	TrustedCAPath          string              `yaml:"trustedCAPath" envconfig:"trusted_ca_path"`
	TrustedCA              []byte
}

// NewConfig returns a Config struct from serialized config file
func NewMultiClusterConfig(configFile string) (*MultiClusterConfig, error) {
	cfg := &MultiClusterConfig{
		HTTPPath:    "",
		ServeTLS:    false,
		CertFile:    "/etc/gangway/tls/tls.crt",
		KeyFile:     "/etc/gangway/tls/tls.key",
		SessionSalt: hardCodedDefaultSalt,
		Host:        "0.0.0.0",
		Port:        8080,
	}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal([]byte(data), cfg)
		if err != nil {
			return nil, err
		}
	}

	err := envconfig.Process("GANGWAY", cfg)
	if err != nil {
		return nil, fmt.Errorf("error processing environment variables for prefix %s: %v", "GANGWAY_", err)
	}

	for env, clusterList := range cfg.Clusters {
		for i := range clusterList {
			cluster := &clusterList[i]
			if cluster.EnvPrefix == "" {
				cluster.EnvPrefix = fmt.Sprintf("%s_CLUSTER%d_", env, i)
			}
			err := envconfig.Process(cluster.EnvPrefix+"GANGWAY", cluster)
			if err != nil {
				return nil, fmt.Errorf("error processing environment variables for prefix %s: %v", cluster.EnvPrefix+"GANGWAY_", err)
			}

			err = cluster.Validate()
			if err != nil {
				return nil, err
			}

			err = loadCerts(cluster, cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	// Check for trailing slash on HTTPPath and remove
	cfg.HTTPPath = strings.TrimRight(cfg.HTTPPath, "/")
	return cfg, nil
}

// Validate verifies all properties of config struct are intialized
func (cfg *Config) Validate() error {
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{cfg.ProviderURL == "", "no providerURL specified"},
		{cfg.ClientID == "", "no clientID specified"},
		{cfg.ClientSecret == "" && !cfg.AllowEmptyClientSecret, "no clientSecret specified"},
		{cfg.RedirectURL == "", "no redirectURL specified"},
		{cfg.APIServerURL == "", "no apiServerURL specified"},
	}

	for _, check := range checks {
		if check.bad {
			return fmt.Errorf("invalid config: %s", check.errMsg)
		}
	}
	return nil
}

// GetRootPathPrefix returns '/' if no prefix is specified, otherwise returns the configured path
func (clusterCfg *MultiClusterConfig) GetRootPathPrefix() string {
	if len(clusterCfg.HTTPPath) == 0 {
		return "/"
	}

	return strings.TrimRight(clusterCfg.HTTPPath, "/")
}

func loadCerts(cfg *Config, clusterCfg *MultiClusterConfig) error {
	if cfg.ClusterCAPath != "" {
		clusterCA, err := os.ReadFile(cfg.ClusterCAPath)
		if err != nil {
			return err
		}
		cfg.ClusterCA = clusterCA
	}

	if clusterCfg.TrustedCAPath != "" {
		trustedCA, err := os.ReadFile(clusterCfg.TrustedCAPath)
		if err != nil {
			return err
		}
		clusterCfg.TrustedCA = trustedCA
	}

	return nil
}
