package router

import (
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/getProperties"
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/getPropertyByID"
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/saveProperty"
	"github.com/alimzhanoff/property-finder/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type storage interface {
	getProperties.AllPropertiesGetter
	getPropertyByID.ByIDGetter
	saveProperty.PropertySaver
}

func InitRoutes(s storage, cfg *config.ServerConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Allow requests from any origin

			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Allow specified HTTP methods

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

			// Allow specified headers

			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

			// Continue with the next handler

			next.ServeHTTP(w, r)

		})
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"), //The url pointing to API definition
	))
	// hello godoc
	// @Summary      Show an account
	// @Description  get string by ID
	// @Tags         accounts
	// @Accept       json
	// @Produce      json
	// @Success      200
	// @Router       /hello [get]
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})
	r.Route("/api/v1/properties", func(r chi.Router) {
		r.Get("/", getProperties.New(s))
		r.Post("/", saveProperty.New(s))

		r.Get("/{property_id}", getPropertyByID.New(s))
	})

	return r
}
