package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/buildsi/drift-analysis/server/database"
)

// removeInflectionPoint removes an inflection point from the database.
func (h *handler) removeInflectionPoint(w http.ResponseWriter, r *http.Request) {
	var point database.InflectionPoint

	err := json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.db.RemovePoint(point)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
