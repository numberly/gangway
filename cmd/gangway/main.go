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
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/numberly/gangway/assets"
	"github.com/numberly/gangway/internal/config"
	"github.com/numberly/gangway/internal/session"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var clusterCfg *config.MultiClusterConfig
var cfg *config.Config
var oauth2Cfg *oauth2.Config

var store *sessions.CookieStore

var sessionManager *session.Session

var transportConfig *config.TransportConfig
var provider *oidc.Provider
var verifier *oidc.IDTokenVerifier

// wrapper function for http logging
func httpLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
		fn(w, r)
	}
}

func rootPathHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if clusterCfg.HTTPPath == "" && r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func main() {
	clusterCfgile := flag.String("config", "", "The config file to use.")
	flag.Parse()

	var err error
	clusterCfg, err = config.NewMultiClusterConfig(*clusterCfgile)
	if err != nil {
		log.Errorf("Could not parse config file: %s", err)
		os.Exit(1)
	}

	transportConfig = config.NewTransportConfig(clusterCfg.TrustedCA)

	var assetFs http.FileSystem
	if clusterCfg.CustomAssetsDir != "" {
		assetFs = http.Dir("clusterCfg.CustomAssetsDir")
	} else {
		assetFs = http.FS(assets.FS)
	}

	store = sessions.NewCookieStore([]byte(clusterCfg.SessionSecurityKey))
	sessionManager = session.New(clusterCfg.SessionSecurityKey, clusterCfg.SessionSalt)

	r := mux.NewRouter()

	r.HandleFunc(clusterCfg.GetRootPathPrefix(), httpLogger(rootPathHandler(clustersHome))).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/login", clusterCfg.HTTPPath), httpLogger(loginHandler)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/callback", clusterCfg.HTTPPath), httpLogger(callbackHandler)).Methods("GET")

	// middleware'd routes
	r.Handle(fmt.Sprintf("%s/logout", clusterCfg.HTTPPath), loginRequired(http.HandlerFunc(logoutHandler))).Methods("GET")
	r.Handle(fmt.Sprintf("%s/commandline", clusterCfg.HTTPPath), loginRequired(http.HandlerFunc(commandlineHandler))).Methods("GET")
	r.Handle(fmt.Sprintf("%s/kubeconf", clusterCfg.HTTPPath), loginRequired(http.HandlerFunc(kubeConfigHandler))).Methods("GET")

	// assets
	assetsPath := fmt.Sprintf("%s/assets/", clusterCfg.HTTPPath)
	r.PathPrefix(assetsPath).Handler(http.StripPrefix(assetsPath, http.FileServer(assetFs)))

	bindAddr := fmt.Sprintf("%s:%d", clusterCfg.Host, clusterCfg.Port)
	// create http server with timeouts
	httpServer := &http.Server{
		Addr:         bindAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if clusterCfg.ServeTLS {
		// update http server with TLS config
		httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12, // minimum TLS 1.2
			// P curve order does not matter, as breaking one means all others can be brute-forced as well:
			// Golang developers prefer:
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384, tls.CurveP521},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_CHACHA20_POLY1305_SHA256, // TLS 1.3
				tls.TLS_AES_256_GCM_SHA384,       // TLS 1.3
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_AES_128_GCM_SHA256, // TLS 1.3
			},
		}
	}

	// start up the http server
	go func() {
		log.Infof("Gangway started! Listening on: %s", bindAddr)
		// exit with FATAL logging why we could not start
		// example: FATA[0000] listen tcp 0.0.0.0:8080: bind: address already in use
		if clusterCfg.ServeTLS {
			log.Fatal(httpServer.ListenAndServeTLS(clusterCfg.CertFile, clusterCfg.KeyFile))
		} else {
			log.Fatal(httpServer.ListenAndServe())
		}
	}()

	// create channel listening for signals so we can have graceful shutdowns
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting.")
	// close the HTTP server
	_ = httpServer.Shutdown(context.Background())
}
