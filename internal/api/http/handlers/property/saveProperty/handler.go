package saveProperty

import (
	"context"
	"errors"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/alimzhanoff/property-finder/internal/storage"
	"github.com/alimzhanoff/property-finder/pkg/api/response"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type PropertySaver interface {
	SavePropertyWithPropertyTypeWithTx(ctx context.Context, property models.Property) error
}

func New(h PropertySaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.saveProperty.New"

		in := &PropertyDTO{}
		if err := render.Bind(r, in); err != nil {
			render.Render(w, r, response.ErrBadRequest)
			return
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(in)
		if err != nil {
			render.Render(w, r, response.ErrBadRequest)
			return
		}

		property := models.Property{
			PropertyType:     &models.PropertyType{TypeName: in.PropertyTypeName},
			Address:          &models.Address{},
			AddressText:      in.AddressText,
			Price:            in.Price,
			Rooms:            &in.Rooms,
			Area:             &in.Area,
			Description:      &in.DescriptionText,
			ConstructionYear: &in.ConstructionYear,
			HasPool:          &in.HasPool,
			DistanceToMetro:  &in.DistanceToMetro,
			MetroStation:     &models.MetroStation{},
		}
		err = h.SavePropertyWithPropertyTypeWithTx(r.Context(), property)
		if err != nil {
			if errors.Is(err, storage.ErrUnexpectedRowCount) {
				render.Render(w, r, response.ErrBadRequest)
				return
			}
			render.Render(w, r, response.ErrInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[int]string{http.StatusOK: "Successfully"})
	}
}
