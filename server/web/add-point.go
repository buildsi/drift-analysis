package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/buildsi/drift-analysis/server/database"
)

//addInflectionPoint handles adding a new inflection point to the database from the API endpoint.
func (h *handler) addInflectionPoint(w http.ResponseWriter, r *http.Request) {
	var point database.InflectionPoint

	err := json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.AddPoint(point)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
