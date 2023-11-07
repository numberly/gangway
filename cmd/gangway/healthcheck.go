package main

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {

	// Créer un ResponseRecorder pour capturer le statut de la réponse
	rec := httptest.NewRecorder()

	// L'appel à clustersHome renvoie maintenant un statut et une erreur.
	status, err := clustersHomeCheck(rec, r)

	// Utiliser le statut et l'erreur renvoyés par clustersHome pour décider de la réponse.
	if status == http.StatusOK && err == nil {
		io.WriteString(w, `{"ready": true}`)
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}
