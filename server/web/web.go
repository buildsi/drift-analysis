package web

import (
	"encoding/json"
	"net/http"

	"github.com/buildsi/drift-analysis/server/database"
	"github.com/buildsi/drift-analysis/server/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type handler struct {
	db *database.DB
	ds datastore.DS
}

func Start(addr string, auth map[string]string, db *database.DB, ds datastore.DS) error {
	// Setup Chi Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	// Set the admin router to the default
	a := r

	// Setup handler functions for api endpoints
	h := handler{
		db: db,
		ds: ds,
	}

	// Enable basic HTTP authentication if the auth map is set.
	if auth != nil {
		a = chi.NewRouter()
		a.Use(middleware.Logger)
		a.Use(middleware.BasicAuth("/", auth))
	}

	// Setup API endpoints
	// No Auth GET Opperations
	r.Get("/artifact/*", h.getArtifact)
	r.Get("/concretizer-diff*", h.getConcretizerDiff)
	r.Get("/inflection-point/*", h.getInflectionPoint)
	r.Get("/inflection-points*", h.getInflectionPoints)
	r.Get("/specs*", h.getSpecs)

	// Auth POST Operations
	// Add Operations
	a.Put("/inflection-point", h.addInflectionPoint)
	a.Put("/artifact*", h.addArtifact)

	// Delete Operations
	a.Delete("/inflection-point", h.removeInflectionPoint)
	a.Delete("/artifact*", h.removeArtifact)

	// Mount admin router
	r.Mount("/", a)

	// Start http server and listen for incoming connections
	return http.ListenAndServe(addr, r)
}

func convertJSON(points []database.InflectionPoint) (string, error) {
	buf, err := json.Marshal(points)
	return string(buf), err
}
