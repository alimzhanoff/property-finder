package saveProperty

import "net/http"

type PropertyDTO struct {
	PropertyTypeName string  `json:"propertyTypeName" validate:"required"`
	AddressText      string  `json:"addressText" validate:"required"`
	Price            float64 `json:"price" validate:"required,min=0"`
	Rooms            int     `json:"rooms,omitempty" validate:"omitempty,min=0"`
	Area             float64 `json:"area,omitempty" validate:"omitempty,min=0"`
	DescriptionText  string  `json:"descriptionText,omitempty"`
	ConstructionYear int     `json:"constructionYear,omitempty" validate:"omitempty,min=0"`
	HasPool          bool    `json:"hasPool,omitempty"`
	DistanceToMetro  float64 `json:"distanceToMetro,omitempty" validate:"omitempty,min=0"`
}

func (p *PropertyDTO) Bind(r *http.Request) error {
	return nil
}
