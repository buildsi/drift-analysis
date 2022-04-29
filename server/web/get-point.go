package web

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/buildsi/drift-analysis/server/database"
)

func (h *handler) getInflectionPoint(w http.ResponseWriter, r *http.Request) {
	var point database.InflectionPoint
	var err error

	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/inflection-point"), "/")
	if len(id) > 0 {
		point, err = h.db.GetPointByID(id)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}

	}

	if point.AbstractSpec == "" {
		http.NotFound(w, r)
		return
	}

	out, err := convertJSON([]database.InflectionPoint{point})
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(15*time.Minute).Format(http.TimeFormat))
	w.Write([]byte(out))
}
