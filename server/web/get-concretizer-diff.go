package web

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/buildsi/drift-analysis/server/database"
)

func (h *handler) getConcretizerDiff(w http.ResponseWriter, r *http.Request) {
	var points []database.InflectionPoint
	var err error

	spec := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/concretizer-diff"), "/")
	if len(spec) > 0 {
		points, err = h.db.GetAllWithSpec(spec)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
	} else {
		points, err = h.db.GetAll()
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
	}

	concretizerMap := make(map[string]map[string]map[string]database.InflectionPoint)
	for _, point := range points {
		if concretizerMap[point.Concretizer] == nil {
			concretizerMap[point.Concretizer] =
				make(map[string]map[string]database.InflectionPoint)
			if concretizerMap[point.Concretizer][point.AbstractSpec] == nil {
				concretizerMap[point.Concretizer][point.AbstractSpec] =
					make(map[string]database.InflectionPoint)
			}
		}
		concretizerMap[point.Concretizer][point.AbstractSpec][point.GitCommit] = point
	}

	result := []database.InflectionPoint{}

	for concretizer := range concretizerMap {
		for comparativeConcretizer := range concretizerMap {
			if concretizer == comparativeConcretizer {
				continue
			}
			for abstractSpec := range concretizerMap[concretizer] {
				for gitCommit, point := range concretizerMap[concretizer][abstractSpec] {
					if _, ok := concretizerMap[comparativeConcretizer][abstractSpec][gitCommit]; ok {
						result = append(result, point)
					}

				}
			}
		}
	}

	out, err := convertJSON(result)
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
