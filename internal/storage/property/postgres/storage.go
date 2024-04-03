package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/alimzhanoff/property-finder/internal/models"
	"github.com/alimzhanoff/property-finder/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	*Postgres
}

const (
	OneRow = 1
	TxKey  = "transaction"
)

func (s *Storage) Store() {

}
func (s *Storage) SavePropertyWithPropertyTypeWithTx(ctx context.Context, property models.Property) error {
	const op = "storage.property.postgres.SavePropertyWithPropertyTypeWithTx"

	var err error = nil

	ctx = context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: get transaction: %v", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	ctx = context.WithValue(ctx, TxKey, tx)

	id, err := s.GetPropertyTypeIDWithTx(ctx, property.PropertyType.TypeName)
	if err != nil {
		return err
	}
	property.PropertyType.ID = id
	if err := s.SavePropertyWithTx(ctx, property); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %v", op, err)
	}
	return nil
}

func (s *Storage) GetPropertyTypeIDWithTx(ctx context.Context, propertyType string) (pt_id int, err error) {
	const op = "storage.property.postgres.GetPropertyTypeID"

	stmt := `
SELECT pt.property_type_id
FROM Property_types as pt
WHERE pt.property_type_name=$1;`

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return 0, fmt.Errorf("%s: cannot get tx from context", op)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = tx.QueryRow(ctx, stmt, propertyType).Scan(&pt_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrNotFound
		}
		return 0, fmt.Errorf("%s: execute statement: %v", op, err)
	}
	return pt_id, nil
}
func (s *Storage) SavePropertyWithTx(ctx context.Context, property models.Property) (err error) {
	const op = "storage.property.postgres.saveProperty"

	stmt := `
INSERT INTO Properties (property_type_id, property_address_text, property_price, property_rooms, property_area, property_description, property_construction_year, property_has_pool, property_distance_to_metro)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("%s: cannot get tx from context", op)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	ct, err := tx.Exec(ctx,
		stmt,
		property.PropertyType.ID,
		property.AddressText,
		property.Price,
		*property.Rooms,
		*property.Area,
		*property.Description,
		*property.ConstructionYear,
		*property.HasPool,
		*property.DistanceToMetro,
	)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %v", op, err)
	}
	r := ct.RowsAffected()
	if r != int64(OneRow) {
		err = fmt.Errorf("%s: %v", op, storage.ErrUnexpectedRowCount)
		return err
	}
	fmt.Println(r)
	return nil
}
func (s *Storage) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	const op = "storage.property.postgres.GetAllProperties"

	properties := make([]models.Property, 0)

	stmt := `SELECT
    p.property_id,
    p.property_type_id,
    pt.property_type_name,
    p.address_id,
    p.property_address_text,
    p.property_price,
    p.property_rooms,
    p.property_area,
    p.property_description,
    p.property_construction_year,
    p.property_has_pool,
    p.property_distance_to_metro,
    p.metro_station_id
FROM Properties AS p
JOIN Property_types AS pt ON p.property_type_id = pt.property_type_id;`

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return properties, fmt.Errorf("%s: execute statement: %v", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		p := models.NewProperty()

		err := rows.Scan(
			&p.ID,
			&p.PropertyType.ID,
			&p.PropertyType.TypeName,
			&p.Address.AddressID,
			&p.AddressText,
			&p.Price,
			&p.Rooms,
			&p.Area,
			&p.Description,
			&p.ConstructionYear,
			&p.HasPool,
			&p.DistanceToMetro,
			&p.MetroStation.ID,
		)

		if err != nil {
			return properties, fmt.Errorf("%s: scan to struct: %v", op, err)
		}

		properties = append(properties, p)
	}

	if err := rows.Err(); err != nil {
		return properties, fmt.Errorf("%s: %v", op, err)
	}

	return properties, nil
}

func (s *Storage) GetPropertyByID(ctx context.Context, id int) (models.Property, error) {
	const op = "storage.property.postgres.GetPropertyByID"

	p := models.NewProperty()

	stmt := `SELECT
    p.property_id,
    p.property_type_id,
    pt.property_type_name,
    p.address_id,
    p.property_address_text,
    p.property_price,
    p.property_rooms,
    p.property_area,
    p.property_description,
    p.property_construction_year,
    p.property_has_pool,
    p.property_distance_to_metro,
    p.metro_station_id
FROM Properties AS p
JOIN Property_types AS pt ON p.property_type_id = pt.property_type_id
WHERE p.property_id=$1;`

	err := s.db.QueryRow(ctx, stmt, id).Scan(
		&p.ID,
		&p.PropertyType.ID,
		&p.PropertyType.TypeName,
		&p.Address.AddressID,
		&p.AddressText,
		&p.Price,
		&p.Rooms,
		&p.Area,
		&p.Description,
		&p.ConstructionYear,
		&p.HasPool,
		&p.DistanceToMetro,
		&p.MetroStation.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, storage.ErrNotFound
		}
		return p, fmt.Errorf("%s: execute statement: %v", op, err)
	}

	return p, nil
}

func New(ctx context.Context, postgres *Postgres) (*Storage, error) {
	const op = "storage.property.postgres.New"
	if postgres == nil {
		return nil, fmt.Errorf("%s: cannot get a storage", op)
	}

	stmt := `
CREATE TABLE IF NOT EXISTS Countries (
                           country_id SERIAL PRIMARY KEY,
                           country_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Regions (
                         region_id SERIAL PRIMARY KEY,
                         region_name VARCHAR(100) NOT NULL,
                         country_id INT NOT NULL,
                         FOREIGN KEY (country_id) REFERENCES Countries(country_id)
);

CREATE TABLE IF NOT EXISTS Cities (
                        city_id SERIAL PRIMARY KEY,
                        city_name VARCHAR(100) NOT NULL,
                        region_id INT NOT NULL,
                        FOREIGN KEY (region_id) REFERENCES Regions(region_id)
);

CREATE TABLE IF NOT EXISTS Streets (
                         street_id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         city_id INT NOT NULL,
                         FOREIGN KEY (city_id) REFERENCES Cities(city_id)
);

CREATE TABLE IF NOT EXISTS Addresses (
                           address_id SERIAL PRIMARY KEY,
                           street_id INT NOT NULL,
                           address_postal_code VARCHAR(20),
                           latitude DECIMAL(10, 8),
                           longitude DECIMAL(11, 8),
                           FOREIGN KEY (street_id) REFERENCES Streets(street_id)
);
CREATE TABLE IF NOT EXISTS Property_types (
                               property_type_id SERIAL PRIMARY KEY,
                               property_type_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS Metro_stations (
                               metro_station_id SERIAL PRIMARY KEY,
                               metro_station_name TEXT
);
CREATE TABLE IF NOT EXISTS Properties (
                          property_id SERIAL PRIMARY KEY,
                          property_type_id INT NOT NULL,
                          address_id INT,
                          property_address_text VARCHAR(100) NOT NULL ,
                          property_price DECIMAL NOT NULL ,
                          property_rooms INTEGER,
                          property_area DECIMAL,
                          property_description TEXT,
                          property_construction_year INTEGER,
                          property_has_pool BOOLEAN,
                          property_distance_to_metro DECIMAL,
                          metro_station_id INT,
                          FOREIGN KEY (property_type_id) REFERENCES Property_types(property_type_id),
                          FOREIGN KEY (address_id) REFERENCES Addresses(address_id),
                          FOREIGN KEY (metro_station_id) REFERENCES Metro_stations(metro_station_id)
);


INSERT INTO Property_types (property_type_name) VALUES ('Квартира'), ('Дом'), ('Земельный участок');

INSERT INTO Countries (country_name) VALUES ('Казахстан');
INSERT INTO Regions (region_name, country_id) VALUES ('ЮКО', 1), ('Алматинская область', 1);
INSERT INTO Cities (city_name, region_id) VALUES ('Шымкент', 1), ('Алматы', 2);
INSERT INTO Streets (name, city_id) VALUES ('Назарбаев', 1), ('Жандосова', 2);
INSERT INTO Addresses (street_id, address_postal_code, latitude, longitude) VALUES (1, '123456', 50.123456, 30.123456), (2, '654321', 40.654321, 20.654321), (2, '654321', 40.654321, 20.654321);
INSERT INTO Metro_stations (metro_station_name) VALUES ('Краснопресненская'),('Баррикадная'),('Пушкинская');

INSERT INTO Properties (property_type_id, address_id, property_address_text, property_price, property_rooms, property_area, property_description, property_construction_year, property_has_pool, property_distance_to_metro, metro_station_id)
VALUES
    (1, 1, 'Текстовый адрес', 100000, 3, 80, 'Описание квартиры 1', 2005, TRUE, 0.5, 1),
    (2, 2, 'Текстовый адрес', 250000, 5, 200, 'Описание дома 1', 2010, TRUE, 1.2, 2),
    (3, 3, 'Текстовый адрес', 250000, 4, 1000, 'Описание участка 1', NULL, FALSE, NULL, 1);`

	_, err := postgres.db.Exec(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{postgres}, nil
}

func Drop(ctx context.Context, postgres *Postgres) error {
	const op = "storage.property.postgres.Drop"
	if postgres == nil {
		return fmt.Errorf("%s: cannot get a storage", op)
	}

	stmt := `
	drop table if exists properties;

	drop table if exists addresses;

	drop table if exists streets;

	drop table if exists cities;

	drop table if exists regions;

	drop table if exists countries;

	drop table if exists property_types;

	drop table if exists metro_stations;`

	_, err := postgres.db.Exec(ctx, stmt)

	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
