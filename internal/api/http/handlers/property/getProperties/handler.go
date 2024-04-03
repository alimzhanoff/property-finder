package getProperties

import (
	"context"
	"errors"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/alimzhanoff/property-finder/internal/storage"
	"github.com/alimzhanoff/property-finder/pkg/api/response"
	"github.com/go-chi/render"
	"net/http"
)

type AllPropertiesGetter interface {
	GetAllProperties(ctx context.Context) ([]models.Property, error)
}

// New godoc
// @Summary      List properties
// @Description  get properties
// @Tags         properties
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Property
// @Failure      404  {object}  response.ErrResponse
// @Failure      500  {object}  response.ErrResponse
// @Router       /properties [get]
func New(h AllPropertiesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.GetAllProperties.New"

		properties, err := h.GetAllProperties(r.Context())
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				render.Render(w, r, response.ErrNotFound)
				return
			}
			render.Render(w, r, response.ErrInternalServerError)
			return
		}

		render.JSON(w, r, properties)
		return
	}
}
