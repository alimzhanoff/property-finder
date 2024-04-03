package models

type Country struct {
	CountryID int    `json:"country_id" db:"country_id"`
	Name      string `json:"country_name" db:"country_name"`
}

type Region struct {
	RegionID int     `json:"region_id" db:"region_id"`
	Name     string  `json:"region_name" db:"region_name"`
	Country  Country `json:"region_country" db:"region_country"`
}

type City struct {
	CityID int    `json:"city_id" db:"city_id"`
	Name   string `json:"city_name" db:"city_name"`
	Region Region `json:"region" db:"region"`
}

type Street struct {
	StreetID int    `json:"street_id" db:"street_id"`
	Name     string `json:"street_name" db:"street_name"`
	City     City   `json:"city" db:"city"`
}

type Address struct {
	AddressID  int     `json:"address_id" db:"address_id"`
	Street     Street  `json:"street" db:"street"`
	PostalCode string  `json:"address_postal_code" db:"address_postal_code"`
	Latitude   float64 `json:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" db:"longitude"`
}

type PropertyType struct {
	ID       int    `json:"property_type_id" db:"property_type_id"`
	TypeName string `json:"property_type_name" db:"property_type_name"`
}

type MetroStation struct {
	ID   int    `json:"metro_station_id" db:"metro_station_id"`
	Name string `json:"metro_station_name" db:"metro_station_name"`
}

type Property struct {
	ID               int           `json:"property_id" db:"property_id"`
	PropertyType     *PropertyType `json:"property_type" db:"property_type"`
	Address          *Address      `json:"property_address,omitempty" db:"property_address,nullable"`
	AddressText      string        `json:"property_address_text" db:"property_address_text"`
	Price            float64       `json:"property_price" db:"property_price"`
	Rooms            *int          `json:"property_rooms,omitempty" db:"property_rooms,nullable"`
	Area             *float64      `json:"property_area,omitempty" db:"property_area,nullable"`
	Description      *string       `json:"property_description,omitempty" db:"property_description,nullable"`
	ConstructionYear *int          `json:"property_construction_year,omitempty" db:"property_construction_year,nullable"`
	HasPool          *bool         `json:"property_has_pool,omitempty" db:"property_has_pool,nullable"`
	DistanceToMetro  *float64      `json:"property_distance_to_metro,omitempty" db:"property_distance_to_metro,nullable"`
	MetroStation     *MetroStation `json:"metro_station,omitempty" db:"metro_station,nullable"`
}

func NewProperty() Property {
	return Property{
		PropertyType: &PropertyType{},
		Address:      &Address{},
		MetroStation: &MetroStation{},
	}
}
