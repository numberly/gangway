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

package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	htmltemplate "html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/numberly/gangway/internal/config"
	"github.com/numberly/gangway/internal/jwt"
	"github.com/numberly/gangway/templates"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

// userInfo stores information about an authenticated user
type userInfo struct {
	ClusterName  string
	Username     string
	Claims       map[string]interface{}
	KubeCfgUser  string
	IDToken      string
	RefreshToken string
	ClientID     string
	ClientSecret string
	IssuerURL    string
	APIServerURL string
	ClusterCA    string
	TrustedCA    string
	ShowClaims   bool
	HTTPPath     string
}

type clusterHomeInfo struct {
	Clusters map[string][]config.Config
	HTTPPath string
}

func serveTemplate(tmplFile string, data interface{}, w http.ResponseWriter) error {
	var (
		templatePath string
		templateData []byte
		err          error
	)

	// Use custom templates if provided
	if clusterCfg.CustomHTMLTemplatesDir != "" {
		templatePath = filepath.Join(clusterCfg.CustomHTMLTemplatesDir, tmplFile)
		templateData, err = os.ReadFile(templatePath)
	} else {
		templateData, err = templates.FS.ReadFile(tmplFile)
	}

	if err != nil {
		log.Errorf("Failed to find template asset: %s at path: %s", tmplFile, templatePath)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	tmpl := htmltemplate.New(tmplFile).Funcs(FuncMap())
	tmpl, err = tmpl.Parse(string(templateData))
	if err != nil {
		log.Errorf("Failed to parse template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tmpl.ExecuteTemplate(w, tmplFile, data)
	if err != nil {
		log.Errorf("Failed to render template %s: %s", tmplFile, err)
	}
	return nil
}

func generateKubeConfig(cfg *userInfo) clientcmdapi.Config {
	// fill out kubeconfig structure
	kcfg := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		CurrentContext: cfg.ClusterName,
		Clusters: []clientcmdapi.NamedCluster{
			{
				Name: cfg.ClusterName,
				Cluster: clientcmdapi.Cluster{
					Server:                   cfg.APIServerURL,
					CertificateAuthorityData: []byte(cfg.ClusterCA),
				},
			},
		},
		Contexts: []clientcmdapi.NamedContext{
			{
				Name: cfg.ClusterName,
				Context: clientcmdapi.Context{
					Cluster:  cfg.ClusterName,
					AuthInfo: cfg.KubeCfgUser,
				},
			},
		},
		AuthInfos: []clientcmdapi.NamedAuthInfo{
			{
				Name: cfg.KubeCfgUser,
				AuthInfo: clientcmdapi.AuthInfo{
					AuthProvider: &clientcmdapi.AuthProviderConfig{
						Name: "oidc",
						Config: map[string]string{
							"client-id":                      cfg.ClientID,
							"client-secret":                  cfg.ClientSecret,
							"id-token":                       cfg.IDToken,
							"idp-issuer-url":                 cfg.IssuerURL,
							"idp-certificate-authority-data": base64.StdEncoding.EncodeToString([]byte(cfg.TrustedCA)),
							"refresh-token":                  cfg.RefreshToken,
						},
					},
				},
			},
		},
	}
	return kcfg
}

func loginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusTemporaryRedirect)
			return
		}

		// Valider le token
		claims, err := jwt.ValidateToken(tokenString, *clusterCfg)
		if err != nil {
			http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clustersHome(w http.ResponseWriter, _ *http.Request) {

	data := &clusterHomeInfo{
		Clusters: clusterCfg.Clusters,
		HTTPPath: clusterCfg.HTTPPath,
	}

	_ = serveTemplate("clustersHome.tmpl", data, w)
}

func clustersHomeCheck(w http.ResponseWriter, _ *http.Request) (int, error) {
	data := &clusterHomeInfo{
		Clusters: clusterCfg.Clusters,
		HTTPPath: clusterCfg.HTTPPath,
	}

	_ = serveTemplate("clustersHome.tmpl", data, w)

	// Serve the template and handle any errors
	if err := serveTemplate("clustersHome.tmpl", data, w); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Handler pour le login.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the cluster name from the request.
	clusterName := r.URL.Query().Get("cluster")
	if clusterName == "" {
		// If no cluster is specified, redirect to the cluster selection page.
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusSeeOther)
		return
	}

	// Obtain the cluster configuration based on the cluster name.
	clusterConfig, ok := getClusterConfig(clusterName)
	if !ok {
		// If the cluster name is not valid, return an error.
		http.Error(w, "Invalid cluster name", http.StatusBadRequest)
		return
	}

	// Create a new OIDC provider using the cluster provider URL.
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, clusterConfig.ProviderURL)
	if err != nil {
		log.Errorf("Could not create OIDC provider for cluster %s: %s", clusterName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create an OIDC verifier to ensure that the received tokens are valid.
	verifier = provider.Verifier(&oidc.Config{ClientID: clusterConfig.ClientID})

	// Configure the OAuth2 client with the cluster information.
	oauth2Cfg = &oauth2.Config{
		ClientID:     clusterConfig.ClientID,
		ClientSecret: clusterConfig.ClientSecret,
		RedirectURL:  clusterConfig.RedirectURL,
		Scopes:       clusterConfig.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	// Generate a random state for the OAuth request and store it in the session.
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		log.Errorf("failed to generate random data: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// Create a JWT with the cluster name and OAuth state.
	signedToken, err := jwt.CreateToken(clusterName, state, *clusterCfg)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the JWT as a secure, HttpOnly cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	})

	// Construct the authentication URL and redirect the client to the OIDC provider.
	authURL := oauth2Cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline, // To request a refresh token.
		oauth2.SetAuthURLParam("prompt", "consent"), // Force the user to give consent.
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusSeeOther)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, transportConfig.HTTPClient)

	// Retrieve the JWT from the cookie.
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "No auth token present", http.StatusUnauthorized)
		return
	}

	// Validate the JWT.
	claims, err := jwt.ValidateToken(cookie.Value, *clusterCfg)
	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusUnauthorized)
		return
	}

	// Check the state for a match.
	state := r.URL.Query().Get("state")
	if state != claims.OAuth2State {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Use the access code to retrieve a token.
	code := r.URL.Query().Get("code")
	oauth2Token, err := oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		log.Errorf("failed to exchange token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID token and verify its validity.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Errorf("no id_token found")
		http.Error(w, "Internal error: no id_token found", http.StatusInternalServerError)
		return
	}

	_, err = verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Errorf("failed to verify token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new JWT with the ID token and refresh token.
	signedToken, err := jwt.UpdateToken(cookie.Value, rawIDToken, oauth2Token.RefreshToken, *clusterCfg)
	if err != nil {
		http.Error(w, "Failed to update auth token", http.StatusInternalServerError)
		return
	}

	// Send the new JWT to the client.
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	})

	// Redirect to the command line page with the cluster name as a parameter.
	http.Redirect(w, r, fmt.Sprintf("%s/commandline?cluster=%s", clusterCfg.HTTPPath, claims.ClusterName), http.StatusSeeOther)
}

func commandlineHandler(w http.ResponseWriter, r *http.Request) {
	info := generateInfo(w, r)
	if info == nil {
		// generateInfo writes to the ResponseWriter if it encounters an error.
		// TODO(abrand): Refactor this.
		return
	}

	_ = serveTemplate("commandline.tmpl", info, w)
}

func kubeConfigHandler(w http.ResponseWriter, r *http.Request) {
	info := generateInfo(w, r)
	if info == nil {
		// generateInfo writes to the ResponseWriter if it encounters an error.
		// TODO(abrand): Refactor this.
		return
	}

	d, err := yaml.Marshal(generateKubeConfig(info))
	if err != nil {
		log.Errorf("Error creating kubeconfig - %s", err.Error())
		http.Error(w, "Error creating kubeconfig", http.StatusInternalServerError)
		return
	}

	// tell the browser the returned content should be downloaded
	w.Header().Add("Content-Disposition", "Attachment")
	_, err = w.Write(d)
	if err != nil {
		log.Errorf("Failed to write kubeconfig: %v", err)
	}
}

func generateInfo(w http.ResponseWriter, r *http.Request) *userInfo {
	// Retrieve the JWT from the cookie.
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "No auth token present", http.StatusUnauthorized)
		return nil
	}

	// Validate the JWT.
	claimsJwt, err := jwt.ValidateToken(cookie.Value, *clusterCfg)
	if err != nil {
		http.Error(w, "Invalid auth token", http.StatusUnauthorized)
		return nil
	}

	rawIDToken := claimsJwt.OAuth2TokenId
	refreshToken := claimsJwt.OAuth2RefreshId

	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, transportConfig.HTTPClient)

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Errorf("failed to verify token: %v", err)
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix()+"/login", http.StatusSeeOther)
		return nil
	}

	claims := make(map[string]interface{})
	if err := idToken.Claims(&claims); err != nil {
		log.Errorf("failed to unmarshal claims: %v", err)
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix()+"/login", http.StatusSeeOther)
		return nil
	}

	clusterName := r.URL.Query().Get("cluster")

	if clusterName == "" {
		// Si aucun cluster n'est spécifié, redirigez vers la page de sélection du cluster.
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusSeeOther)
		return nil
	}

	cfg, ok := getClusterConfig(clusterName)
	if !ok {
		http.Error(w, "Invalid cluster name", http.StatusBadRequest)
		return nil
	}

	username, ok := claims[cfg.UsernameClaim].(string)
	if !ok {
		http.Error(w, "Could not parse Username claim", http.StatusInternalServerError)
		return nil
	}

	kubeCfgUser := strings.Join([]string{username, cfg.ClusterName}, "@")

	issuerURL, ok := claims["iss"].(string)
	if !ok {
		http.Error(w, "Could not parse Issuer URL claim", http.StatusInternalServerError)
		return nil
	}

	if cfg.ClientSecret == "" {
		log.Warn("Setting an empty Client Secret should only be done if you have no other option and is an inherent security risk.")
	}

	info := &userInfo{
		ClusterName:  cfg.ClusterName,
		Username:     username,
		Claims:       claims,
		KubeCfgUser:  kubeCfgUser,
		IDToken:      rawIDToken,
		RefreshToken: refreshToken,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		IssuerURL:    issuerURL,
		APIServerURL: cfg.APIServerURL,
		ClusterCA:    string(cfg.ClusterCA),
		TrustedCA:    string(clusterCfg.TrustedCA),
		ShowClaims:   cfg.ShowClaims,
		HTTPPath:     clusterCfg.HTTPPath,
	}
	return info
}

func getClusterConfig(clusterName string) (config.Config, bool) {
	for _, clusters := range clusterCfg.Clusters {
		for _, cluster := range clusters {
			if cluster.ClusterName == clusterName {
				return cluster, true
			}
		}
	}
	return config.Config{}, false
}
