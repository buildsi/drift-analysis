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

type Spec struct {
	Name                string
	NumDependencies     int
	NumInflectionPoints int
}

type JsonSpec struct {
	Nodes []JsonNodes `json:"nodes"`
}

type JsonNodes struct {
	Dependencies []interface{} `json:"dependencies"`
}

func (h *handler) getSpecs(w http.ResponseWriter, r *http.Request) {
	var points []database.InflectionPoint
	var err error

	// Get inflection points for spec or all specs.
	spec := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/specs"), "/")
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

	// Trim points to just the latest inflection point for
	// each spec.
	results := []Spec{}
	latestPoints := make(map[string]database.InflectionPoint)
	numPoints := make(map[string]int)

	for _, ip := range points {
		numPoints[ip.AbstractSpec]++
		if latestPoints[ip.AbstractSpec].GitCommitDate.Before(ip.GitCommitDate) &&
			ip.Concretized {
			latestPoints[ip.AbstractSpec] = ip
		}
	}

	for abstractSpec, point := range latestPoints {
		concreteSpecJson, err := h.ds.Get(point.SpecUUID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
		concreteSpecData := make(map[string]JsonSpec)
		err = json.Unmarshal([]byte(concreteSpecJson), &concreteSpecData)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", 500)
			return
		}
		numDeps := len(concreteSpecData["spec"].Nodes) - 1

		results = append(results, Spec{
			Name:                abstractSpec,
			NumInflectionPoints: numPoints[abstractSpec],
			NumDependencies:     numDeps,
		})
	}

	out, err := json.Marshal(results)
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
