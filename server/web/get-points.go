package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/buildsi/drift-analysis/server/database"
)

func (h *handler) getInflectionPoints(w http.ResponseWriter, r *http.Request) {
	var out string

	spec := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/inflection-points"), "/")
	if len(spec) > 0 {
		points, err := h.db.GetAllWithSpec(spec)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
		out, err = convertJSON(points)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
	} else {
		points, err := h.db.GetAll()
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
		out, err = convertJSON(points)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(15*time.Minute).Format(http.TimeFormat))
	w.Write([]byte(out))
}

func convertJSON(points []database.InflectionPoint) (string, error) {
	buf, err := json.Marshal(points)
	return string(buf), err
}
