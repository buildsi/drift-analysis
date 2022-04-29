package web

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// getArtifact retrieves an artifact from the datastore given its UUID.
func (h *handler) getArtifact(w http.ResponseWriter, r *http.Request) {
	// Parse the uuid from the input url.
	uuid := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/artifact"), "/")

	// Get the artifact from the database.
	payload, err := h.ds.Get(uuid)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	artifactType, err := h.db.GetArtifactType(uuid)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	// Write out the payload.
	w.Header().Set("Content-Type", artifactType)
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(15*time.Minute).Format(http.TimeFormat))
	w.Write([]byte(payload))
}
