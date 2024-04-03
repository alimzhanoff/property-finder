package getPropertyByID_test

import (
	"fmt"
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/getPropertyByID"
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/getPropertyByID/mocks"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newint(num int) *int {
	return &num
}
func newstr(s string) *string {
	return &s
}
func newfloat(num float64) *float64 {
	return &num
}
func newbool(b bool) *bool {
	return &b
}

func TestPropertyByIDHandler(t *testing.T) {
	cases := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "[{\"property_id\":1,\"property_type\":{\"property_type_id\":1,\"property_type_name\":\"Квартира\"},\"property_address_text\":\"some address\",\"property_price\":100,\"property_rooms\":3,\"property_area\":300,\"property_description\":\"some description\",\"property_construction_year\":2000,\"property_has_pool\":false,\"property_distance_to_metro\":3000}]\n",
		},
		{
			name:           "error - invalid property ID",
			id:             "aws",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Bad Request"}`,
		},
		{
			name:           "error - property not found",
			id:             "2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Not Found"}`,
		},
		{
			name:           "error - internal server error",
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error"}`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks.ByIDGetter{}

			m.On("GetPropertyByID", mock.Anything, tt.id).Return(models.Property{
				ID:               1,
				PropertyType:     &models.PropertyType{ID: 1, TypeName: "Квартира"},
				AddressText:      "some address",
				Price:            100,
				Rooms:            newint(3),
				Area:             newfloat(300),
				Description:      newstr("some description"),
				ConstructionYear: newint(2000),
				HasPool:          newbool(false),
				DistanceToMetro:  newfloat(3000),
			}, nil)

			r := chi.NewRouter()
			r.Get("/api/v1/properties/{property_id}", getPropertyByID.New(m))

			ts := httptest.NewServer(r)
			defer ts.Close()
			handler := getPropertyByID.New(m)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/properties/%s", tt.id), nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		})
	}
}
