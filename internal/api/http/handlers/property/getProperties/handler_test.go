package getProperties_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/alimzhanoff/property-finder/internal/api/http/handlers/property/getProperties"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/alimzhanoff/property-finder/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:generate mockery --name=AllPropertiesGetter --output=mocks --outpkg=mocks
type mockAllPropertiesGetter struct {
	mock.Mock
}

func (m *mockAllPropertiesGetter) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Property), args.Error(1)
}

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
func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockGetter := &mockAllPropertiesGetter{}
		mockGetter.On("GetAllProperties", mock.Anything).Return([]models.Property{
			{1, &models.PropertyType{1, "Квартира"}, nil, "some address", 100, newint(3), newfloat(300), newstr("some description"), newint(2000), newbool(false), newfloat(3000), nil},
		}, nil)

		handler := getProperties.New(mockGetter)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/properties", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "[{\"property_id\":1,\"property_type\":{\"property_type_id\":1,\"property_type_name\":\"Квартира\"},\"property_address_text\":\"some address\",\"property_price\":100,\"property_rooms\":3,\"property_area\":300,\"property_description\":\"some description\",\"property_construction_year\":2000,\"property_has_pool\":false,\"property_distance_to_metro\":3000}]\n", rr.Body.String())
		mockGetter.AssertExpectations(t)
	})

	t.Run("error - not found", func(t *testing.T) {
		mockGetter := &mockAllPropertiesGetter{}
		mockGetter.On("GetAllProperties", mock.Anything).Return([]models.Property{}, storage.ErrNotFound)

		handler := getProperties.New(mockGetter)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/properties", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		fmt.Println(rr.Body.String())

		// Проверяем, что код ответа соответствует ожидаемому
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Проверяем, что тело ответа содержит правильное сообщение об ошибке
		expectedBody := "{\"status\":\"Resource not found.\"}\n"
		assert.Equal(t, expectedBody, rr.Body.String())

		// Проверяем, что ожидания mock-объекта были выполнены
		mockGetter.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		mockGetter := &mockAllPropertiesGetter{}
		mockGetter.On("GetAllProperties", mock.Anything).Return([]models.Property{}, errors.New("some error"))

		handler := getProperties.New(mockGetter)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/properties", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		fmt.Println(rr.Body.String())

		// Проверяем, что код ответа соответствует ожидаемому
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Проверяем, что тело ответа содержит правильное сообщение об ошибке
		expectedBody := "{\"status\":\"Internal Server Error\"}\n"
		assert.Equal(t, expectedBody, rr.Body.String())

		// Проверяем, что ожидания mock-объекта были выполнены
		mockGetter.AssertExpectations(t)
	})
}
