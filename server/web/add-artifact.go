package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// addArtifact adds an artifact to the datastore and returns its uuid
func (h *handler) addArtifact(w http.ResponseWriter, r *http.Request) {
	// Parse the uuid from the input url.
	id := uuid.New()
	for _, err := h.ds.Get(id.String()); err == nil; {
		id = uuid.New()
	}

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ds.Put(id.String(), buf.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the artifact type to the database
	err = h.db.AddArtifactType(id.String(), r.Header["Content-Type"][0])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write out the payload.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(15*time.Minute).Format(http.TimeFormat))
	w.Write([]byte(fmt.Sprintf("{\"uuid\":\"%s\"}\n", id.String())))
}
