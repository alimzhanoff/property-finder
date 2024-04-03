CREATE TABLE IF NOT EXISTS Countries (
                           country_id SERIAL PRIMARY KEY,
                           name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Regions (
                         region_id SERIAL PRIMARY KEY,
                         name VARCHAR(100) NOT NULL,
                         country_id INT NOT NULL,
                         FOREIGN KEY (country_id) REFERENCES Countries(country_id)
);

CREATE TABLE IF NOT EXISTS Cities (
                        city_id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
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
                           postal_code VARCHAR(20),
                           latitude DECIMAL(10, 8),
                           longitude DECIMAL(11, 8),
                           FOREIGN KEY (street_id) REFERENCES Streets(street_id)
);
CREATE TABLE IF NOT EXISTS Property_types (
                               property_type_id SERIAL PRIMARY KEY,
                               type_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS Properties (
                          property_id SERIAL PRIMARY KEY,
                          property_type_id INTEGER REFERENCES Property_types(id),
                          address_id INT REFERENCES Addresses(address_id),
                          address_text VARCHAR(100),
                          price DECIMAL NOT NULL ,
                          rooms INTEGER NOT NULL,
                          area DECIMAL NOT NULL,
                          description TEXT NOT NULL,
                          construction_year INTEGER,
                          has_pool BOOLEAN,
                          distance_to_metro DECIMAL,
                          metro_station_id INTEGER REFERENCES Metro_stations(id)
);
CREATE TABLE IF NOT EXISTS Metro_stations (
                               metro_stations_id SERIAL PRIMARY KEY,
                               name TEXT
);


INSERT INTO Property_types (type_name) VALUES ('Квартира'), ('Дом'), ('Земельный участок');

INSERT INTO Countries (name) VALUES ('Казахстан');
INSERT INTO Regions (name, country_id) VALUES ('ЮКО', 1), ('Алматинская область', 1);
INSERT INTO Cities (name, region_id) VALUES ('Шымкент', 1), ('Алматы', 2);
INSERT INTO Streets (name, city_id) VALUES ('Назарбаев', 1), ('Жандосова', 2);
INSERT INTO Addresses (street_id, postal_code, latitude, longitude) VALUES (1, '123456', 50.123456, 30.123456), (2, '654321', 40.654321, 20.654321);

INSERT INTO Properties (property_type_id, address_id, address_text, price, rooms, area, description, construction_year, has_pool, distance_to_metro, metro_station_id)
VALUES
    (1, 1, 'Текстовый адрес', 100000, 3, 80, 'Описание квартиры 1', 2005, TRUE, 0.5, 1),
    (2, 2, 'Текстовый адрес', 250000, 5, 200, 'Описание дома 1', 2010, TRUE, 1.2, 2),
    (3, 3, 'Текстовый адрес', 250000, 4, 1000, 'Описание участка 1', NULL, FALSE, NULL, NULL);

CREATE OR REPLACE FUNCTION update_address_text()
RETURNS TRIGGER AS $$
BEGIN
UPDATE Properties
SET address_text = CONCAT(
        (SELECT Streets.name FROM Streets WHERE Streets.street_id = NEW.address_id),
        ', ',
        (SELECT Cities.name FROM Cities JOIN Streets ON Cities.city_id = Streets.city_id WHERE Streets.street_id = NEW.address_id),
        ', ',
        (SELECT Regions.name FROM Regions JOIN Cities ON Regions.region_id = Cities.region_id JOIN Streets ON Cities.city_id = Streets.city_id WHERE Streets.street_id = NEW.address_id),
        ', ',
        (SELECT Countries.name FROM Countries JOIN Regions ON Countries.country_id = Regions.country_id JOIN Cities ON Regions.region_id = Cities.region_id JOIN Streets ON Cities.city_id = Streets.city_id WHERE Streets.street_id = NEW.address_id)
                   )
WHERE id = NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_address_text_trigger
    BEFORE INSERT ON Properties
    FOR EACH ROW
    EXECUTE FUNCTION update_address_text();
