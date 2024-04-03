package getPropertyByID

import (
	"context"
	"errors"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/alimzhanoff/property-finder/internal/storage"
	"github.com/alimzhanoff/property-finder/pkg/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

//go:generate mockery --name=ByIDGetter --output=mocks --outpkg=mocks
type ByIDGetter interface {
	GetPropertyByID(ctx context.Context, id int) (models.Property, error)
}

// New godoc
// @Summary      Show a property
// @Description  get property by ID
// @Tags         property
// @Accept       json
// @Produce      json
// @Param        property_id   path      int  true  "property_id"
// @Success      200  {object}  models.Property
// @Failure      404  {object}  response.ErrResponse
// @Failure      500  {object}  response.ErrResponse
// @Router       /properties/{property_id} [get]
func New(h ByIDGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.getPropertyByID"

		propertyIDStr := chi.URLParam(r, "property_id")
		propertyID, err := strconv.Atoi(propertyIDStr)
		if err != nil {
			render.Render(w, r, response.ErrNotFound)
			return
		}

		property, err := h.GetPropertyByID(context.Background(), propertyID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				render.Render(w, r, response.ErrNotFound)
				return
			}
			render.Render(w, r, response.ErrInternalServerError)
			return
		}

		render.JSON(w, r, property)
	}
}

//func New(h ByIDGetter) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		const op = "http.handlers.getPropertyByID"
//
//		propertyIDStr := chi.URLParam(r, "property_id")
//		propertyID, err := strconv.Atoi(propertyIDStr)
//		fmt.Println("HIIIIIIIIIIIIIIIIIIIII: ", propertyIDStr)
//		if err != nil {
//			render.Render(w, r, response.ErrBadRequest)
//			return
//		}
//
//		property, err := h.GetPropertyByID(context.Background(), propertyID)
//		if err != nil {
//			if errors.Is(err, storage.ErrNotFound) {
//				render.Render(w, r, response.ErrNotFound)
//				return
//			}
//			render.Render(w, r, response.ErrInternalServerError)
//			return
//		}
//
//		render.JSON(w, r, property)
//	}
//}
