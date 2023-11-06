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

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/numberly/gangway/internal/config"
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

func serveTemplate(tmplFile string, data interface{}, w http.ResponseWriter) {
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
		return
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
		gangwayIdTokenSess, err := sessionManager.Session.Get(r, "gangway_id_token")
		if err != nil {
			log.Errorf("Unable to get session : %v", err)
			http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusTemporaryRedirect)
			return
		}

		if gangwayIdTokenSess.Values["id_token"] == nil {
			log.Error("id_token is nil")
			http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clustersHome(w http.ResponseWriter, _ *http.Request) {

	data := &clusterHomeInfo{
		Clusters: clusterCfg.Clusters,
		HTTPPath: clusterCfg.HTTPPath,
	}

	serveTemplate("home.tmpl", data, w)
}

// Handler pour le login.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le nom du cluster à partir de la requête.
	clusterName := r.URL.Query().Get("cluster")
	if clusterName == "" {
		// Si aucun cluster n'est spécifié, redirigez vers la page de sélection du cluster.
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusSeeOther)
		return
	}

	// Obtenez la configuration du cluster en fonction du nom du cluster.
	clusterConfig, ok := getClusterConfig(clusterName)
	if !ok {
		// Si le nom du cluster n'est pas valide, renvoyez une erreur.
		http.Error(w, "Invalid cluster name", http.StatusBadRequest)
		return
	}

	// Créez un nouveau fournisseur OIDC en utilisant l'URL du fournisseur du cluster.
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, clusterConfig.ProviderURL)
	if err != nil {
		log.Errorf("Could not create OIDC provider for cluster %s: %s", clusterName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Créer un vérificateur OIDC pour s'assurer que les tokens reçus sont valides.
	verifier = provider.Verifier(&oidc.Config{ClientID: clusterConfig.ClientID})

	// Configurer le client OAuth2 avec les informations du cluster.
	oauth2Cfg = &oauth2.Config{
		ClientID:     clusterConfig.ClientID,
		ClientSecret: clusterConfig.ClientSecret,
		RedirectURL:  clusterConfig.RedirectURL,
		Scopes:       clusterConfig.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	// Générer un état aléatoire pour la requête OAuth et le stocker dans la session.
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		log.Errorf("failed to generate random data: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// Utiliser initSession pour initialiser la session.
	gangwaySess, err := sessionManager.Session.Get(r, "gangway")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gangwaySess.Values["state"] = state
	gangwaySess.Values["clusterName"] = clusterName
	err = gangwaySess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Construit l'URL d'authentification et redirige le client vers le fournisseur OIDC.
	authURL := oauth2Cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline, // Pour demander un token d'actualisation.
		oauth2.SetAuthURLParam("prompt", "consent"), // Forcer l'utilisateur à donner son consentement.
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	sessionManager.Cleanup(w, r, "gangway")
	sessionManager.Cleanup(w, r, "gangway_id_token")
	sessionManager.Cleanup(w, r, "gangway_refresh_token")

	// Redirection après logout.
	http.Redirect(w, r, clusterCfg.GetRootPathPrefix(), http.StatusSeeOther)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, transportConfig.HTTPClient)

	// Charger la session principale.
	session, err := sessionManager.Session.Get(r, "gangway")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier le nom du cluster dans la session.
	clusterName, ok := session.Values["clusterName"].(string)
	if !ok || clusterName == "" {
		http.Error(w, "Internal error: clusterName not found", http.StatusInternalServerError)
		return
	}

	// Vérifier la concordance de l'état.
	state := r.URL.Query().Get("state")
	if state != session.Values["state"] {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Utiliser le code d'accès pour récupérer un token.
	code := r.URL.Query().Get("code")
	oauth2Token, err := oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		log.Errorf("failed to exchange token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extraire le ID token et vérifier sa validité.
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

	gangwayIdTokenSess, err := sessionManager.Session.Get(r, "gangway_id_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gangwayIdTokenSess.Values["id_token"] = rawIDToken
	err = gangwayIdTokenSess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	gangwayRefreshTokenSess, err := sessionManager.Session.Get(r, "gangway_refresh_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gangwayRefreshTokenSess.Values["refresh_token"] = oauth2Token.RefreshToken
	err = gangwayRefreshTokenSess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Rediriger vers la page de ligne de commande avec le nom du cluster comme paramètre.
	http.Redirect(w, r, fmt.Sprintf("%s/commandline?cluster=%s", clusterCfg.HTTPPath, clusterName), http.StatusSeeOther)
}

func commandlineHandler(w http.ResponseWriter, r *http.Request) {
	info := generateInfo(w, r)
	if info == nil {
		// generateInfo writes to the ResponseWriter if it encounters an error.
		// TODO(abrand): Refactor this.
		return
	}

	serveTemplate("commandline.tmpl", info, w)
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
	// Load the session.
	sessionIdToken, err := sessionManager.Session.Get(r, "gangway_id_token")
	if err != nil {
		log.Errorf("Error retrieving session: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil
	}

	rawIDToken, ok := sessionIdToken.Values["id_token"].(string)
	if !ok {
		log.Errorf("id_token is not OK : %v", ok)
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix()+"/login", http.StatusSeeOther)
		return nil
	}

	sessionRefreshToken, err := sessionManager.Session.Get(r, "gangway_refresh_token")
	if err != nil {
		log.Errorf("Error retrieving session: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil
	}

	refreshToken, ok := sessionRefreshToken.Values["refresh_token"].(string)
	if !ok {
		log.Errorf("refresh_token is not OK : %v", ok)
		http.Redirect(w, r, clusterCfg.GetRootPathPrefix()+"/login", http.StatusSeeOther)
		return nil
	}

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
