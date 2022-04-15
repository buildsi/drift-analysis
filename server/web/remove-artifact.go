package web

import (
	"log"
	"net/http"
	"strings"
)

// removes an artifact from the datastore given its UUID.
func (h *handler) removeArtifact(w http.ResponseWriter, r *http.Request) {
	// Parse the uuid from the input url.
	uuid := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/artifact"), "/")

	// Get the artifact from the database.
	err := h.ds.Delete(uuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
